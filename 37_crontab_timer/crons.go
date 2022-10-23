package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc("0 0 0 10 * *", func() { fmt.Println("每10秒") })
	c.AddFunc("@hourly", func() { fmt.Println("每小时") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("每小时三十") })
	c.Start()
	// ..
	// Funcs 在它们自己的 goroutine 中异步调用。
	// ...
	// Funcs 也可以添加到正在运行的 Cron
	c.AddFunc("@daily", func() { fmt.Println("每天") })
	// ..
	// 检查 cron 作业条目的下一次和上一次运行时间。
	for ind, ent := range c.Entries() {
		fmt.Printf("index:%v entries:%v\n", ind, ent)
	}

	// ..
	c.Stop() // 停止调度程序（不停止任何已经运行的作业）。
}
