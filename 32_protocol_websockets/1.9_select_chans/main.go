package main

import (
	"fmt"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
)

func select_chans() {

	var c1 = make(chan int, 1)
	var c2 = make(chan float64, 1)
	var c3 = make(chan string, 2)
	wg.Add(3)
	go func() {
		c1 <- 10
		wg.Done()
	}()

	go func() {
		c2 <- 3.1415926
		wg.Done()
	}()

	go func() {
		c3 <- "Hello. "
		c3 <- "World"
		wg.Done()
	}()

	for i := 0; i < 5; i++ {
		//// 如果不遍历，可能随机取一个，
		/// 多少个chan 就遍历多数次
		/// 如果有多余的 遍历次数，将执行default
		select {
		case cc1 := <-c1:
			fmt.Println("cc1 get:", cc1)
		case cc2 := <-c2:
			fmt.Println("cc2 get:", cc2)
		case cc3 := <-c3:
			fmt.Println("cc3 get:", cc3)
		default:
			fmt.Println("default.")
		}
	}
	wg.Wait()
}

type mystring struct {
	Value string
}

func MakeNewStr(ms string) *mystring {
	var mss mystring

	ms = fmt.Sprintf("%v.%v", "Hello", ms)
	mss.Value = ms
	return &mss
}

type New_methods interface {
	GetUpper()
	GetLower()
}

func (ms *mystring) GetUpper() string {
	return strings.ToUpper(ms.Value)
}
func (ms *mystring) GetLower() string {
	return strings.ToLower(ms.Value)
}

func type_interface() {

}
func main() {
	// select_chans()

	ms1 := MakeNewStr("hello.")
	ms2 := MakeNewStr("world.")
	fmt.Println(ms1.GetLower(), ms1.GetUpper())
	fmt.Println(ms2.GetLower(), ms2.GetUpper())
}
