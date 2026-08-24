package main

import "fmt"

func Filter[T any](in <-chan T, filter func(T) bool) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for value := range in {
			if filter(value) {
				result <- value
			}
		}
	}()

	return result
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)

		for i := 0; i < 10; i++ {
			channel <- i
		}
	}()

	isOdd := func(value int) bool {
		return value%2 != 0
	}

	for value := range Filter(channel, isOdd) {
		fmt.Println(value)
	}
}
