package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4}
	b := a[2:3] // b = [3]
	b = append(b, 7)
	fmt.Println(a, len(a), cap(a)) // [1 2 3 7] 4 4
	fmt.Println(b, len(b), cap(b)) // [3 7] 2 2
}
