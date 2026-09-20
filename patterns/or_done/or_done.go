package main

import (
	"fmt"
	"time"
)

func OrDone[T any](in chan T, done chan struct{}) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			// приоритезация, т.к. в select рандомный выбор
			select {
			case <-done:
				return
			default:
			}

			select {
			case value, opened := <-in:
				if !opened {
					return
				}
				out <- value
			case <-done:
				return
			}
		}
	}()

	return out
}

func main() {
	channel := make(chan string)

	go func() {
		for {
			channel <- "test"
			time.Sleep(200 * time.Millisecond)
		}
	}()

	done := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()

	/*
		for {
			select{
			case value, opened := <- in:
				if !opened{
					return
				}
				// processing
			case <- done:
				return
			}
		}
	*/

	for value := range OrDone(channel, done) {
		fmt.Println(value)
	}
}
