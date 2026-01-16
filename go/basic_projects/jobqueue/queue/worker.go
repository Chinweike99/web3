package queue

import (
	"context"
	"fmt"
	"time"

	"jobqueue/job"
)

// StartWorker launches a worker goroutine
func (q *Queue) StartWorker(ctx context.Context, id int) {
	q.wg.Add(1)

	go func() {
		defer q.wg.Done()
		fmt.Printf("Worker %d started\n", id)

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d draining remaining jobs\n", id)
				q.drain(ctx)
				return

			case j := <-q.jobs:
				q.handleJob(ctx, id, j)
			}
		}
	}()
}

// handleJob executes a job safely with timeout and panic recovery
func (q *Queue) handleJob(ctx context.Context, workerID int, j job.Job) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Worker %d recovered from panic: %v\n", workerID, r)
			q.failed <- j
		}
	}()

	jobCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	fmt.Printf("Worker %d processing job %d\n", workerID, j.ID)

	if err := j.Execute(jobCtx); err != nil {
		fmt.Printf("Worker %d job %d failed: %v\n", workerID, j.ID, err)
		q.failed <- j
	}
}

// drain processes remaining jobs before shutdown
func (q *Queue) drain(ctx context.Context) {
	for {
		select {
		case j := <-q.jobs:
			_ = j.Execute(ctx)
		default:
			return
		}
	}
}
