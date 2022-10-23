package main

// 实验channel 传递channel
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	rch <-chan string
	wch chan<- string
	ach chan string
	as  string
)

func handle(wg *sync.WaitGroup, a int) chan int {
	out := make(chan int)
	go func() {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		out <- a
		wg.Done()
	}()
	return out
}

func main() {
	// 只读，只写channel
	// myc := make(chan string, 10)
	// rch = myc
	// wch = myc
	// ach = myc
	// wch <- "a"
	// fmt.Println("len wch", len(wch))
	// fmt.Println("len ach", len(ach))
	// ach <- "b"
	// fmt.Println("len ach", len(ach), len(rch), len(wch))
	// myc <- "c"
	// <-rch
	// fmt.Println("len ach", len(rch))

	// 使用channel 传递channel
	reqs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	outs := make(chan chan int, len(reqs))
	var wg sync.WaitGroup
	wg.Add(len(reqs))
	for _, x := range reqs {
		k := handle(&wg, x)
		outs <- k
	}
	go func() {
		wg.Wait()
		close(outs)
	}()
	//读取结果，结果有序
	for k := range outs {
		fmt.Println(<-k)
	}
}
