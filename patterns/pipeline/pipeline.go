package main

import "fmt"

func generate[T any](values ...T) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for _, value := range values {
			result <- value
		}
	}()

	return result
}

func process[T any](in <-chan T, action func(T) T) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for value := range in {
			result <- action(value)
		}
	}()

	return result
}

func main() {
	values := []int{1, 2, 3, 4, 5}
	mul := func(value int) int {
		return value * value
	}

	for value := range process(generate(values...), mul) {
		fmt.Println(value)
	}
}
