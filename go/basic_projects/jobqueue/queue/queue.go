package queue

import (
	"jobqueue/job"
	"sync"
)


type Queue struct {
	jobs chan job.Job
	wg sync.WaitGroup
}


func NewQueue(bufferSize int) *Queue {
	return &Queue{
		jobs: make(chan job.Job, bufferSize),
	}
}

func (q *Queue) Submit(j job.Job) {
	q.jobs <-j
}

func (q *Queue) Wait(){
	q.wg.Wait()
}	

