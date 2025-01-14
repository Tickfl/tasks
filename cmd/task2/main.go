package main

import (
	"flag"
	"fmt"
	"patterns/internal/semaphore"
	"sync"
	"time"
)

func main() {
	var T, N, M int
	flag.IntVar(&T, "T", 3, "Общее количество задач")
	flag.IntVar(&N, "N", 10, "Общее число потоков")
	flag.IntVar(&M, "M", 4, "Максимальное кол-во потоков на одну задачу")
	flag.Parse()

	globalSemaphore := semaphore.NewSemaphore(N)

	var wg sync.WaitGroup

	for i := 0; i < T; i++ {
		wg.Add(1)

		taskId := i + 1
		taskChannel := make(chan int)

		go func(taskId int, taskChannel chan int) {
			defer close(taskChannel)
			for j := 1; j < 10; j++ {
				taskChannel <- j * taskId
			}
		}(taskId, taskChannel)

		go func(taskId int, taskChannel chan int) {
			defer wg.Done()

			taskSemaphore := semaphore.NewSemaphore(M)
			var taskWg sync.WaitGroup

			for data := range taskChannel {
				taskWg.Add(1)

				globalSemaphore.Acquire()
				taskSemaphore.Acquire()

				go func(taskId, data int) {
					defer taskWg.Done()
					defer globalSemaphore.Release()
					defer taskSemaphore.Release()

					fmt.Printf("Task %d: data %d\n", taskId, data)
					time.Sleep(time.Second * 10)
				}(taskId, data)
			}

			taskWg.Wait()
		}(taskId, taskChannel)
	}

	wg.Wait()
	fmt.Println("The End")
}
