package queue

import (
	"errors"
	"sync"

	"jobqueue/job"
)

type Queue struct {
	jobs   chan job.Job
	failed chan job.Job
	wg     sync.WaitGroup
}

// NewQueue creates a bounded job queue
func NewQueue(bufferSize int) *Queue {
	return &Queue{
		jobs:   make(chan job.Job, bufferSize),
		failed: make(chan job.Job, bufferSize),
	}
}

// Submit attempts to enqueue a job without blocking
func (q *Queue) Submit(j job.Job) error {
	select {
	case q.jobs <- j:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// Wait blocks until all workers exit
func (q *Queue) Wait() {
	q.wg.Wait()
}
