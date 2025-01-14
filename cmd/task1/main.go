package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	var once = sync.Once{}

	for _, c := range channels {
		go func(ch <-chan interface{}) {
			for {
				select {
				case _, ok := <-ch:
					if !ok {
						once.Do(func() {
							close(done)
						})
					}
				case <-done:
					return
				}
			}
		}(c)
	}

	return done
}
