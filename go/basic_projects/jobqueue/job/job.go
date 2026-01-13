package job

import "fmt"

type Job struct {
	ID int
	Name string
}

func (j Job) Execute(){
	fmt.Printf("Executing job %d: %s\n", j.ID, j.Name)
}


