package job

import (
	"context"
	"fmt"
	"time"
)

type Job struct {
	ID   int
	Name string
}

func (j Job) Execute(ctx context.Context) error {
	if j.Name == "panic" {
		panic("job panicked intentionally")
	}

	select {
	case <-time.After(3 * time.Second):
		fmt.Printf("Job %d completed: %s\n", j.ID, j.Name)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
