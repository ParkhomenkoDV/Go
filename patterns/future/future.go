package main

import (
	"fmt"
	"time"
)

type Future[T any] struct {
	result chan T
}

func NewFuture[T any](action func() T) Future[T] {
	future := Future[T]{
		result: make(chan T),
	}

	go func() {
		defer close(future.result)
		future.result <- action()
	}()

	return future
}

func (f *Future[T]) Get() T {
	return <-f.result
}

func main() {
	job := func() any {
		time.Sleep(3 * time.Second)
		return "success"
	}

	future := NewFuture(job)
	result := future.Get()
	fmt.Println(result)
}
