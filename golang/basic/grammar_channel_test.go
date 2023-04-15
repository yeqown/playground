package basic_test

import (
	"fmt"
	"testing"
	"time"
)

type Worker struct {
	Id       string
	CostSecs int
}

func (w *Worker) Work(jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker: ", w.Id, "is dealing with job --- ", j)
		time.Sleep(time.Duration(w.CostSecs) * time.Second)
		results <- j * 2
	}
}

func Test_channel(t *testing.T) {
	jobs := make(chan int, 10)
	results := make(chan int, 20)
	workers := make([]Worker, 3)
	for w := 0; w < 3; w++ {
		workers[w] = Worker{
			fmt.Sprintf("worker-%d", w+1),
			w + 1,
		}
		go workers[w].Work(jobs, results)
	}

	// send jobs
	for j := 0; j < 10; j++ {
		jobs <- j
	}
	close(jobs)

	// get result
	for r := 0; r < 9; r++ {
		<-results
	}
}
