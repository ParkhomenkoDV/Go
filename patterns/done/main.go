package main

import "fmt"

func process(done <-chan struct{}) <-chan struct{} {
	out := make(chan struct{})

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			default:
				// processing
			}
		}
	}()

	return out
}

func main() {
	done := make(chan struct{})
	result := process(done)

	close(done)
	<-result

	fmt.Println("terminated")
}
