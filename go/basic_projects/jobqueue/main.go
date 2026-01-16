package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"jobqueue/job"
	"jobqueue/queue"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewQueue(2)

	// Start workers
	for i := 1; i <= 3; i++ {
		q.StartWorker(ctx, i)
	}

	// OS signal handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nShutdown signal received")
		cancel()
	}()

	scanner := bufio.NewScanner(os.Stdin)
	jobID := 1

	fmt.Println("Enter jobs (type 'exit' to quit):")

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "exit" {
			cancel()
			break
		}

		j := job.Job{
			ID:   jobID,
			Name: text,
		}

		if err := q.Submit(j); err != nil {
			fmt.Println("Submit failed:", err)
		} else {
			jobID++
		}
	}

	q.Wait()
	fmt.Println("All workers stopped. Exiting.")
}
