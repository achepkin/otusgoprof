package service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {

	if maxErrorsCount <= 0 {
		maxErrorsCount = len(tasks) + 1
	}

	if len(tasks) < workersCount {
		workersCount = len(tasks)
	}

	tt := make(chan Task, len(tasks))
	result := make(chan error, maxErrorsCount)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	var errCount int32

	for i := 0; i < workersCount; i++ {
		//fmt.Println("worker:", i)
		wg.Add(1)
		go func(i int) {
			fmt.Println("goroutine:", i)
			defer wg.Done()

			for {
				select {
				case <-done:
					fmt.Println("exit gr:", i)
					return
				default:
					// do nothing
				}

				select {
				case task := <-tt:
					result <- task()
				default:
					fmt.Println("sleep:", i)
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(i)
	}

	for _, task := range tasks {
		tt <- task
	}

	go func() {
		doneTaskCount := 0
		var once sync.Once
		for err := range result {
			doneTaskCount++
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			}
			if maxErrorsCount == int(atomic.LoadInt32(&errCount)) || len(tasks) == doneTaskCount {
				once.Do(func() {
					close(done)
				})
			}
			fmt.Println("task done:", doneTaskCount)
		}
	}()

	wg.Wait()
	close(result)

	if int(atomic.LoadInt32(&errCount)) >= maxErrorsCount {
		return ErrErrorsLimitExceeded
	}

	return nil
}
