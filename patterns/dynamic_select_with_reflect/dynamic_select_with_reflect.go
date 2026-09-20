package main

import (
	"fmt"
	"reflect"
)

func main() {
	channel := make(chan int, 1)
	vch := reflect.ValueOf(channel)

	succeed := vch.TrySend(reflect.ValueOf(100))
	fmt.Println(succeed, vch.Len(), vch.Cap())

	branches := []reflect.SelectCase{
		{Dir: reflect.SelectDefault},
		{Dir: reflect.SelectRecv, Chan: vch},
	}

	index, value, ok := reflect.Select(branches)
	fmt.Println(index, value, ok)
}
