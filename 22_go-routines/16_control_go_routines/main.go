package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

var (
	task_cnt = math.MaxInt64
	wg       = sync.WaitGroup{}
	chs      = NewChanCtl() //make(chan int)
	RWMux    = sync.RWMutex{}
)

//不控制数量，启动最大 9223372036854775807 个协程
func RountineNumBug() {

	//max 9223372036854775807
	fmt.Println("math max int64:", task_cnt)
	for i := 0; i < task_cnt; i++ {
		go func(i int) {
			//显示当前有多少 go 协程
			fmt.Println("go func", i, "goroutine count =", runtime.NumGoroutine())
		}(i)
	}
}

//使用缓冲channel 控制go 协程创建数量，和同步模块wg 控制执行数量
func routChan(ch chan bool, i int) {
	fmt.Printf("go func:%v, go routine count:%v\n", i, runtime.NumGoroutine())
	<-ch //channel决定函数的结束时间
	wg.Done()
}

func ChanControlRoutine() {
	task_cnt = 10
	ch := make(chan bool, 3)
	//for 循环决定routChan的创建速度，创建完成后，并不能等到协程全部执行完成
	for i := 0; i < task_cnt; i++ {
		wg.Add(1)
		ch <- true
		go routChan(ch, i)
		fmt.Printf("start go routine cnt:%v\n", i)
		// wg.Wait() //在协程等待执行完成
	}
	wg.Wait() // 在主线程等待 协程执行完成
}

///通过无缓存channel 控制协程数量,读写分离 加锁
type ChanCtrl struct {
	Read  <-chan int
	Write chan<- int
	Locks sync.RWMutex
}

func NewChanCtl() *ChanCtrl {
	ch := make(chan int)
	return &ChanCtrl{
		Read:  ch,
		Write: ch,
		Locks: sync.RWMutex{},
	}
}

func (c *ChanCtrl) Putin(i int) {
	c.Locks.Lock()
	defer c.Locks.Unlock()
	c.Write <- i
}
func (c *ChanCtrl) ReadOne() int {
	c.Locks.Lock()
	defer c.Locks.Unlock()
	return <-c.Read
}

func NoBuffChan() {
	for t := range chs.Read {
		fmt.Println("go task =", t, ", goroutine count =", runtime.NumGoroutine())
		wg.Done()
	}
}

func sendTask(task int) {
	wg.Add(1)
	chs.Putin(task)

}

func NoBuffTaskChan() {

	goCnt := 3
	for i := 0; i < goCnt; i++ {
		//启动协程
		go NoBuffChan()
	}

	taskCnt := math.MaxInt64 //模拟用户需求数量
	for t := 0; t < taskCnt; t++ {
		//执行发送
		sendTask(t)
	}
	wg.Wait()
}

///直接使用 chan 读写分离，不加锁
var chUnlock = make(chan int)

func NoBuffChanUnLock() {
	for t := range chUnlock {
		fmt.Println("go task =", t, ", goroutine count =", runtime.NumGoroutine())
		wg.Done()
	}
}

func sendTaskUnlock(task int) {
	wg.Add(1)
	chUnlock <- task

}

func NoBuffTaskChanUnlock() {

	goCnt := 3
	for i := 0; i < goCnt; i++ {
		//启动协程
		go NoBuffChanUnLock()
	}

	taskCnt := math.MaxInt64 //模拟用户需求数量
	for t := 0; t < taskCnt; t++ {
		//执行发送
		sendTaskUnlock(t)
	}
	wg.Wait()
}

func main() {
	// 不控制go协程 数量系统IDE 编译器 将占满内存
	// RountineNumBug()

	//通过缓冲channel控制
	// wgs := sync.WaitGroup{}
	// wgs.Add(1)
	// go func() {
	// 	ChanControlRoutine()
	// 	wgs.Done()
	// }()
	// wgs.Wait()

	//通过无缓冲 chan控制,加锁顺序执行
	// NoBuffTaskChan()

	//无缓冲通道 不加锁 NoBuffTaskChanUnlock
	NoBuffTaskChanUnlock()
}
