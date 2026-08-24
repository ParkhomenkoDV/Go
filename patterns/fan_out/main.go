package main

import (
	"fmt"
	"sync"
)

func FanOut[T any](in <-chan T, splits int) []<-chan T {
	outs := make([]chan T, splits)
	for i := 0; i < splits; i++ {
		outs[i] = make(chan T)
	}

	go func() {
		idx := 0
		for value := range in {
			select {
			case outs[idx] <- value:
			default:
			}

			idx = (idx + 1) % splits // if idx == splits { idx = 0 }
		}

		for _, ch := range outs {
			close(ch)
		}
	}()

	// can not cast []chan T to []<-chan T
	results := make([]<-chan T, splits)
	for i := 0; i < splits; i++ {
		results[i] = outs[i]
	}

	return results
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)

		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()

	splits := 2

	channels := FanOut(channel, splits)

	var wg sync.WaitGroup
	wg.Add(splits)

	go func() {
		defer wg.Done()

		for value := range channels[0] {
			fmt.Println("ch1: ", value)
		}
	}()

	go func() {
		defer wg.Done()

		for value := range channels[1] {
			fmt.Println("ch2: ", value)
		}
	}()

	wg.Wait()
}
