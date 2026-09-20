package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(nTickets uint) Semaphore {
	return Semaphore{
		tickets: make(chan struct{}, nTickets),
	}
}

func (s *Semaphore) Acquare() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	semaphore := NewSemaphore(5)
	for i := 0; i < 6; i++ {
		semaphore.Acquare()

		go func() {
			defer func() {
				wg.Done()
				semaphore.Release()
			}()

			fmt.Println("working")
			time.Sleep(2 * time.Second)
			fmt.Println("Existing")
		}()
	}

	wg.Wait()
}
