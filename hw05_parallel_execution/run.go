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

	tasksChannel := make(chan Task)
	stopChannel := make(chan struct{})
	var once sync.Once

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, &errorCounter, m, stopChannel, tasksChannel, &once)
	}

	shouldStop := false //иначе из внешнего цикла не выйти (на моем уровне знаний)
	for _, t := range tasks {
		select {
		case <-stopChannel:
			shouldStop = true
		case tasksChannel <- t:
		}

		if shouldStop {
			break
		}
	}

	close(tasksChannel)

	wg.Wait() //важно ставить в конце, чтобы не дать main завершиться

	if errorCounter.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// channel уже содержит указатели внутри себя на реальные данные в памяти -> передаем по значению
// WaitGroup, Once - нет -> передаем по указателю
func worker(wg *sync.WaitGroup, errorCounter *atomic.Int64, m int, stopChannel chan struct{}, tasksChannel chan Task, once *sync.Once) {
	defer wg.Done()
	for {
		select {
		case <-stopChannel:
			return
		case task, ok := <-tasksChannel: //ok значит только то, что в канале были данные и их вычитали
			if !ok {
				return
			}

			if err := task(); err != nil { //err - уже результат выполнения таски
				newErrCounter := errorCounter.Add(1)
				if newErrCounter >= int64(m) {
					once.Do(func() {
						close(stopChannel)
					})
					return
				}
			}
		}
	}
}
