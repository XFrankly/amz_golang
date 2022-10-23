package main

import (
	"fmt"
)

func random(n int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := 0; i < n; i++ {
			select {
			case c <- 0:
			case c <- 1:
			}
		}
	}()
	return c
}

func main() {
	var ints = []int{}
	for i := range random(100) {
		fmt.Println(i)
		ints = append(ints, i)
	}
}
