package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/*
select chan 的工作方式测试:
	异步启动 10个线程做一个需要不同时间执行的事情
	select 选择两个通道的数据
	如果有数据则取出，如果没有则阻塞
	default 可以定义，但不做具体的事情
*/
var (
	wg            = sync.WaitGroup{} //只能使用一个waitgroup
	SigCount  int = 0
	Ncr           = NewMyCount(1)
	Muxrw         = sync.RWMutex{}
	ChClose       = make(chan int)
	TotalLoop int = 10
	LongStory string
)

type Schan string
type Ichan int
type MyCount struct {
	Read    <-chan Schan
	Write   chan<- Schan
	Total   chan int //存总值
	maxsize int
}

func NewMyCount(size int) *MyCount {
	c := make(chan Schan, size)
	t := make(chan int, 1)
	return &MyCount{
		Read:    c,
		Write:   c,
		Total:   t,
		maxsize: size,
	}
}

type BewInterface interface {
	PutIn(i Schan) bool
	LookCount() int
}

type CountAdapter struct {
	BewInterface
}

func NewBewinters(size int) *CountAdapter {
	return &CountAdapter{NewMyCount(size)}
}

func (that *MyCount) PutIn(i Schan) bool {
	Muxrw.Lock()
	defer Muxrw.Unlock()
	if len(that.Write) < 50 {
		that.Write <- i
		SigCount += 1
		if len(that.Total) == 0 {

			that.Total <- SigCount
		}

		return true
	} else {
		return false
	}
}

func (that *MyCount) ReadOne() (Schan, bool) {
	Muxrw.Lock()
	defer Muxrw.Unlock()
	if len(that.Read) > 0 {

		return <-that.Read, true
	} else {
		return "", false
	}
}
func (that *MyCount) LookCount() int {
	Muxrw.Lock()
	defer Muxrw.Unlock()
	if len(that.Total) > 0 {

		<-that.Total
	}
	that.Total <- SigCount
	return SigCount
}

func doThing(n int, name int) int {
	var s int
	for i := 0; i < TotalLoop; i++ {
		Ncr.PutIn(Schan(fmt.Sprintf("thread-%v num:%v", name, i)))
		t := 300 * n
		time.Sleep(time.Millisecond * time.Duration(t))
		s += i
	}

	//最后一个线程 需要执行的时间最长，并且完成了执行循环，可以判断退出条件
	if name == TotalLoop {
		ChClose <- 1
	}
	wg.Done()

	return s

}
func GoFuncEx(n int, name int) {

	wg.Add(1)
	go func() {
		doThing(n, name)
	}()

	wg.Wait()

}

func main() {
	fmt.Printf("equal? %v \n", "rtsp://127.0.0.1/demo" == "rtsp://127.0.0.1/demo")

	for j := 0; j <= TotalLoop; j++ {
		//启动
		fmt.Println("start at thread:", j)
		go func() {
			GoFuncEx(j+1, j)
		}()

	}
	x := 0
	for {
		select {
		case sn := <-Ncr.Read:
			fmt.Printf("get str from ch:%v default case length:%v\n", sn, x)
			x += 1
		case close := <-ChClose:
			fmt.Println("total t:", x)
			if close == 1 {
				os.Exit(close)
			}
			// default:
			// x += 1
			// fmt.Println("read one fail. data chanel length:", len(Ncr.Read), x)
			// time.Sleep(time.Millisecond * 100)
		}
	}

}
