package main

import "fmt"

func Bridge[T any](in chan chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for ch := range in {
			for value := range ch { // maybe lock! Do new gorutine!
				out <- value
			}
		}
	}()

	return out
}

func main() {
	channelChannels := make(chan chan string)

	go func() {
		channel1 := make(chan string, 3)
		for i := 0; i < 3; i++ {
			channel1 <- "1"
		}
		close(channel1)

		channel2 := make(chan string, 3)
		for i := 0; i < 3; i++ {
			channel2 <- "2"
		}
		close(channel2)

		channelChannels <- channel1
		channelChannels <- channel2
		close(channelChannels)
	}()

	for value := range Bridge(channelChannels) {
		fmt.Println(value)
	}
}
