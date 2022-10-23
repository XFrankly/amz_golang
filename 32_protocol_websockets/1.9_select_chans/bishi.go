package main

import (
	"fmt"
	"os"
	"sync"
)

// func Bird() {}
// 输出 "Gua"
// func Dog() {}
// 输出"Wang"
// func Cat() {}
// 输出 "Miao"

//  分别开启三个协程
//  循环输出 gua  ，wang ， miao ，
// 在 元音字母 a, u, i ,输出总数达到50 个以后，停止输出。

type Sp struct {
	Write chan<- string
	Read  <-chan string
	Size  int
	Total int
}

var (
	s     = make(chan string)
	close = make(chan int)
	SSP   = NewSp()
	Mux   = sync.RWMutex{}
	wg    = sync.WaitGroup{}
)

func NewSp() *Sp {
	ss := make(chan string)
	return &Sp{Write: ss, Read: ss}
}

func (that *Sp) Putins(sp string) bool {
	if len(that.Read) < that.Size {

		Mux.Lock()
		defer Mux.Unlock()
		that.Write <- sp
		return true
	}
	return false
}

func (that *Sp) ReadOne() string {
	if len(that.Read) > 0 {

		Mux.Lock()
		defer Mux.Unlock()

		return <-that.Read
	}
	return ""
}

func Parse() {
	num := 1
	for {
		select {
		case str := <-s:
			Counts(&num, str)
		}
	}
}

func Counts(num *int, str string) {
	for _, o := range str {
		switch o {
		case 'a':
			if *num >= 50 {
				// close <- 1
				break
			}
			fmt.Printf("字符串:%s,第%d 个元音字母,元音:%s\n", str, *num, string(o))
			*num++
			SSP.Total++
		case 'u':
			if *num >= 50 {
				// close <- 1
				break
			}
			fmt.Printf("字符串:%s,第%d 个元音字母,元音:%s\n", str, *num, string(o))
			*num++
			SSP.Total++
		case 'i':
			if *num >= 50 {
				// close <- 1
				break
			}
			fmt.Printf("字符串:%s,第%d 个元音字母,元音:%s\n", str, *num, string(o))
			*num++
			SSP.Total++
		}
	}
}

func Birds() {
	// SSP.Putins("Gua")
	Mux.Lock()
	defer Mux.Unlock()
	s <- "Gua"
}
func Dogs() {
	// SSP.Putins("Wang")
	Mux.Lock()
	defer Mux.Unlock()
	s <- "Wang"
}
func Cats() {
	// SSP.Putins("Miao")
	Mux.Lock()
	defer Mux.Unlock()
	s <- "Miao"
}
func main() {

	// wg.Add(2)

	go Parse()
	// for i := 0; i < 50; i++ {
	go func() {
		for {
			Birds()
			Dogs()
			Cats()
		}
	}()
	// }

	select {
	case <-close:
		fmt.Println(" 退出")
		os.Exit(1)
	}

	// wg.Wait()
}
