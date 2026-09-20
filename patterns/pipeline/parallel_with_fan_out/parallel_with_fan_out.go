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

func parse(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		for data := range in {
			out <- fmt.Sprintf("parsed - %s", data)
		}
	}()

	return out
}

func send(in <-chan string, workers int) <-chan string {
	var wg sync.WaitGroup
	wg.Add(workers)

	splitted := FanOut(in, workers)
	out := make(chan string)

	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()

			for data := range splitted[idx] {
				out <- fmt.Sprintf("sent - %s", data)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	channel := make(chan string)

	go func() {
		defer close(channel)

		for i := 0; i < 5; i++ {
			channel <- "value"
		}
	}()

	for value := range send(parse(channel), 2) {
		fmt.Println(value)
	}
}
