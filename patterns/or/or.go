package main

import (
	"fmt"
	"time"
)

func Or[T any](chans ...<-chan T) <-chan T {
	switch len(chans) {
	case 0:
		return nil
	case 1:
		return chans[0]
	}

	out := make(chan T)

	go func() {
		defer close(out)

		switch len(chans) {
		case 2:
			select {
			case <-chans[0]:
			case <-chans[1]:
			}
		default:
			select {
			case <-chans[0]:
			case <-chans[1]:
			case <-chans[2]:
			case <-Or(chans[3:]...):
			}
		}
	}()

	return out
}

func main() {
	start := time.Now()

	<-Or(
		time.After(2*time.Hour),
		time.After(5*time.Minute),
		time.After(1*time.Second),
		time.After(1*time.Hour),
		time.After(10*time.Second),
	)

	fmt.Printf("Called after: %s", time.Since(start))
}
