package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	var wg sync.WaitGroup
	var errorCounter atomic.Int64

	for i := 0; i < len(tasks); i++ {
		if errorCounter.Load() <= int64(m) {
			go worker(tasks[i], &wg, &errorCounter, m)
		}
	}

	wg.Wait() //важно ставить в конце, чтобы не дать main завершиться

	if errorCounter.Load() > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(task Task, wg *sync.WaitGroup, errorCounter *atomic.Int64, m int) {
	//if errorCounter.Load() <= int64(m) {
	wg.Add(1)
	result := task()
	if result != nil {
		errorCounter.Add(1)
	}
	wg.Done()
	//}
}
