package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/otusgoprof/hw05_parallel_execution/service"
)

func main() {
	tasksCount := 50
	tasks := make([]service.Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 23
	err := service.Run(tasks, workersCount, maxErrorsCount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("runTasksCount:", runTasksCount)
}
