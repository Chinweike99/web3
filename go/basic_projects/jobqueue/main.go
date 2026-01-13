package main

import (
	"bufio"
	"context"
	"fmt"
	"jobqueue/job"
	"jobqueue/queue"
	"os"
	"os/signal"
	"strings"
	"syscall"
)


func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewQueue(5)

	for i := 1; i <= 3; i++ {
		q.StartWorker(ctx, i)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nShutting down ....")
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

		j := job.Job {
			ID: jobID,
			Name: text,
		}
		q.Submit(j)
		jobID++
	}
	q.Wait()
	fmt.Println("All workers stopped")

}

