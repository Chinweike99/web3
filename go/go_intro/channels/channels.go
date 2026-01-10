package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LogEvent struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}



type LogJob struct {
	Event  LogEvent
	Result chan error
}

func logWorker(id int, jobs <-chan LogJob) {
	for job := range jobs {
		fmt.Printf(
			"Worker %d processing log: [%s] %s\n",
			id,
			job.Event.Level,
			job.Event.Message,
		)

		time.Sleep(500 * time.Millisecond)

		job.Result <- nil
	}
}

func startWorkers(n int, jobs <-chan LogJob) {
	for i := 1; i <= n; i++ {
		go logWorker(i, jobs)
	}
}

type LogHandler struct {
	jobs chan LogJob
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event LogEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if event.Level == "" || event.Message == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	result := make(chan error, 1)

	job := LogJob{
		Event:  event,
		Result: result,
	}

	select {
	case h.jobs <- job:
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	case <-time.After(1 * time.Second):
		http.Error(w, "Service busy", http.StatusGatewayTimeout)
		return
	}

	select {
	case err := <-result:
		if err != nil {
			http.Error(w, "Processing failed", http.StatusInternalServerError)
			return
		}
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Log received")
}

// ---------- Main ----------

func main() {
	jobQueue := make(chan LogJob, 10)

	startWorkers(3, jobQueue)

	handler := &LogHandler{
		jobs: jobQueue,
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	fmt.Println("Server running on http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}

	shutdown(context.Background(), server)
}

// ---------- Graceful Shutdown ----------

func shutdown(ctx context.Context, server *http.Server) {
	fmt.Println("Shutting down server...")
	server.Shutdown(ctx)
}
