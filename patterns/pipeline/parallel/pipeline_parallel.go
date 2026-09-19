package main

import (
	"fmt"
	"sync"
)

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

	out := make(chan string)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for data := range in {
				out <- fmt.Sprintf("sent - %s", data)
			}
		}()
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
