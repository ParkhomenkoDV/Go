package main

import (
	"log"
	"sync"
)

type Barrier struct {
	mutex sync.Mutex
	count uint
	size  uint

	before chan struct{}
	after  chan struct{}
}

func NewBarrier(size uint) *Barrier {
	return &Barrier{
		size:   size,
		before: make(chan struct{}, size),
		after:  make(chan struct{}, size),
	}
}

func (b *Barrier) Before() {
	b.mutex.Lock()

	b.count++
	if b.count == b.size { // если заполнили нужное количество задач
		for i := 0; i < int(b.size); i++ {
			b.before <- struct{}{} // пишем все вместе
		}
	}

	b.mutex.Unlock()
	<-b.before // Либо Lock, либо читка после записи всех весте
}

func (b *Barrier) After() {
	b.mutex.Lock()

	b.count--
	if b.count == 0 {
		for i := 0; i < int(b.size); i++ {
			b.after <- struct{}{}
		}
	}

	b.mutex.Unlock()
	<-b.after
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	bootstrap := func() {
		log.Println("bootstrap")
	}

	work := func() {
		log.Println("work")
	}

	var count uint = 3
	barrier := NewBarrier(count)
	for i := 0; i < int(count); i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < int(count); j++ {
				// wait for all workersto finish previous loop
				barrier.Before()
				bootstrap()
				// wait for other workers to bootstrap
				barrier.After()
				work()
			}
		}()
	}

	wg.Wait()
}
