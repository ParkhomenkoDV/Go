package main

import (
	"fmt"
	"time"
)

type result[T any] struct {
	value T
	err   error
}

type Promise[T any] struct {
	result chan result[T]
}

func NewPromise[T any](fn func() (T, error)) Promise[T] {
	promise := Promise[T]{
		result: make(chan result[T]),
	}

	go func() {
		defer close(promise.result)

		value, err := fn()
		promise.result <- result[T]{value: value, err: err}
	}()

	return promise
}

func (p *Promise[T]) Then(successFn func(T), errorFn func(error)) {
	go func() {
		result := <-p.result
		if result.err == nil {
			successFn(result.value)
		} else {
			errorFn(result.err)
		}
	}()
}

func main() {
	job := func() (string, error) {
		time.Sleep(1 * time.Second)
		return "ok", nil
	}

	promise := NewPromise(job)
	promise.Then(
		func(value string) {
			fmt.Println("success", value)
		},
		func(err error) {
			fmt.Println("error", err.Error())
		},
	)

	time.Sleep(2 * time.Second)
}
