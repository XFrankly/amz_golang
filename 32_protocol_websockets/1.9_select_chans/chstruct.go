package main

// func Bird() {}
// 输出 "Gua"
// func Dog() {}
// 输出"Wang"
// func Cat() {}
// 输出 "Miao"

//  分别开启三个协程
//  循环输出 gua  ，wang ， miao ，
// 在 元音字母 a, u, i ,输出总数达到50 个以后，停止输出。

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

//测试
var (
	tt        = time.Now()
	Count     = 0
	SigCount  = 0 //单通道计数
	ChanCount = 0 //多通道计数
	vc        = make(chan int, 50)
	vc1       = make(chan int, 50)
	vc2       = make(chan int, 50)
	cchan     = NewBewinters(50)
	wgs       sync.WaitGroup
	Muxrw     sync.RWMutex
	Logs      = log.New(os.Stderr, "DEBUG -", 18)
)

type MyCount struct {
	Read    <-chan int
	Write   chan<- int
	Total   chan int //存总值
	maxsize int
}

func NewMyCount(size int) *MyCount {
	c := make(chan int, size)
	t := make(chan int, 1)
	return &MyCount{
		Read:    c,
		Write:   c,
		Total:   t,
		maxsize: size,
	}
}

type BewInterface interface {
	PutIn(i int) bool
	LookCount() int
}

type CountAdapter struct {
	BewInterface
}

func NewBewinters(size int) *CountAdapter {
	return &CountAdapter{NewMyCount(size)}
}

func (that *MyCount) PutIn(i int) bool {
	Mux.Lock()
	defer Mux.Unlock()
	if len(that.Write) < 50 {
		that.Write <- i
		SigCount += i
		if len(that.Total) == 0 {

			that.Total <- SigCount
		}

		return true
	} else {
		return false
	}
}

func (that *MyCount) ReadOne() (int, bool) {
	Mux.Lock()
	defer Mux.Unlock()
	if len(that.Read) > 0 {

		return <-that.Read, true
	} else {
		return 0, false
	}
}
func (that *MyCount) LookCount() int {
	Mux.Lock()
	defer Mux.Unlock()
	if len(that.Total) > 0 {

		<-that.Total
	}
	that.Total <- SigCount
	return SigCount
}

func Number(s string) int {
	var t int
	t += strings.Count(s, "a")
	t += strings.Count(s, "u")
	t += strings.Count(s, "i")
	return t
}

//变量计数
func NumberByCount(s string) {
	if Count >= 50 {
		return
	}
	t := Number(s)
	Count += t

}

//通道计数
func NumberByChannelStruct(s string) {
	//先查看总数 达到50则不继续
	rst := cchan.LookCount()
	if rst >= 50 {

		// os.Exit(1)

		return
	}

	t := Number(s)
	cchan.PutIn(t)

}

//控制使用哪个方式
func NumbersCtl(s string, ch int) {

	if ch == 1 {

		//使用Channel 速度快
		NumberByChannelStruct(s)
	} else if ch == 0 {

		//全局变量计数器 相比chennel稍慢
		NumberByCount(s)
	} else {

		// 使用多个通道时， 计数器空转
		return
	}
	return
}

func Bird(s string, ch int) {
	// Logs.Println(s)
	NumbersCtl(s, ch)
}
func Dog(s string, ch int) {
	// Logs.Println(s)
	NumbersCtl(s, ch)
}
func Cat(s string, ch int) {
	// Logs.Println(s)
	NumbersCtl(s, ch)
}

//使用通道计数器
func DoIt(ch int) {
	for i := 0; i < 50; i++ {
		if ch == 0 {
			if Count >= 50 {
				Logs.Println("变量计数:循环次数", i, Count)

				break
			}
		}
		if ch == 1 {
			if SigCount >= 50 {
				Logs.Println("单通道计数:循环次数", i, SigCount)

				break
			}
		}

		wgs.Add(3)
		go func() {
			Bird("gua", ch)
			wgs.Done()
		}()
		go func() {
			Dog("wang", ch)
			wgs.Done()
		}()
		go func() {
			Cat("miao", ch)
			wgs.Done()
		}()
		wgs.Wait()

	}
}

//多个通道
func DoItByMulChannels(ch int) {
	for i := 0; i < 50; i++ {
		if ChanCount >= 50 {
			Logs.Println("3通道计数:循环次数", i, ChanCount)

			break
		}
		wgs.Add(3)
		go func() {
			words := "gua"
			Bird(words, ch)
			vc <- Number(words)
			wgs.Done()
		}()
		go func() {
			words := "wang"
			Dog(words, ch)
			vc1 <- Number(words)
			wgs.Done()
		}()
		go func() {
			words := "miao"
			Cat(words, ch)
			vc2 <- Number(words)
			wgs.Done()
		}()
		wgs.Wait()

		select {
		case cc1 := <-vc:
			ChanCount += cc1
		case cc2 := <-vc1:
			ChanCount += cc2
		case cc3 := <-vc2:
			ChanCount += cc3
		default:
			Logs.Println("Multi Channel default.")
		}

	}

}

func main() {

	//不用channel 使用变量计数器 75，77，55毫秒
	DoIt(0)

	//使用单个channel 82，55，46  毫秒
	DoIt(1)

	//使用多个channel 63，47，74毫秒
	DoItByMulChannels(2)

	Logs.Println("总耗时", time.Since(tt).Milliseconds(), "毫秒")

}
