package transformer

import "fmt"

func Transformer[T any](in <-chan T, action func(T) T) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for v := range in {
			result <- action(v)
		}
	}()

	return result
}

func main() {
	channel := make(chan int)

	go func() {
		defer close(channel)

		for i := 0; i < 5; i++ {
			channel <- i
		}
	}()

	mul := func(value int) int {
		return value * value
	}

	for number := range Transformer(channel, mul) {
		fmt.Println(number)
	}
}
