package main

import (
	"fmt"
	"time"
)

func foo() {
	fmt.Println("start sleep.")
	time.Sleep(time.Second * 1)
	fmt.Println("End sleep.")
}
func TimeTicker() {
	ticker := time.NewTicker(time.Second * 2)
	// 清理计时器
	defer ticker.Stop()
	for i := 0; i <= 10; i++ {
		fmt.Println("ticker start", time.Now().Format("15:04:05"))
		foo()
		<-ticker.C // ticker 自动新建了一个 channel
	}
}

func TimeTimer() {
	//t *testing.T
	timer := time.NewTimer(time.Second * 2)
	// 清理计时器
	defer timer.Stop()

	for j := 0; j <= 5; j++ {
		fmt.Println("timer start. No.", j, time.Now().Format("15:04:05"))
		foo()
		<-timer.C      // CPU执行后 将销毁管道，只执行一次
		timer.Reset(1) // time.Duration int64
	}
}

func main() {
	TimeTicker()
	TimeTimer()
}
