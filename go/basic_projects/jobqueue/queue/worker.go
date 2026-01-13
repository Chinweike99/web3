package queue

import (
	"context"
	"fmt"
)



func (q *Queue) StartWorker(ctx context.Context, id int) {
	q.wg.Add(1)

	go func ()  {
		defer q.wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d stopping\n", id)
				return

			case j := <-q.jobs:
				fmt.Printf("Worker %d picked job %d\n", id, j.ID)
				j.Execute()
			}
		}
	}()
}

