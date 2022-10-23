package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func OnceDoExample(n int) {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only Once.")
	}
	done := make(chan bool)
	for i := 0; i < n; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < n; i++ {
		x := <-done
		fmt.Println("get from done:", x)
	}
}

//////////////////////池的使用
var bufPool = sync.Pool{
	New: func() any {
		// 池的 New 函数 应该仅生成用于 返回 指针类型
		// Pool 的 New 函数通常应该只返回指针类型，因为可以将指针放入返回接口类型中而无需分配
		return new(bytes.Buffer)
	},
}

/// timeNow 是一个假的用于测试的 time.Now 版本
func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}
func LogPoolDoExample(w io.Writer, key, val string) {
	/*
			池是一组可以单独保存和检索的临时对象。
		存储在池中的任何项目都可能随时自动删除，恕不另行通知。
		如果在发生这种情况时 Pool 拥有唯一的引用，则该项目可能会被释放。

		一个 Pool 可以安全地同时被多个 goroutine 使用。

		Pool 的目的是缓存已分配但未使用的项目以供以后重用，减轻垃圾收集器的压力。
		也就是说，它使构建高效、线程安全的空闲列表变得容易。
		但是，它并不适用于所有空闲列表。

		Pool 的适当用途是管理一组在包的并发独立客户端之间静默共享并可能被重用的临时项目。
		Pool 提供了一种在许多客户端之间分摊分配开销的方法。

		一个很好地使用池的例子是在 fmt 包中，它维护一个动态大小的临时输出缓冲区存储。
		存储在负载下扩展（当许多 goroutine 正在积极打印时）并在静止时缩小。

		另一方面，作为短期对象的一部分维护的空闲列表不适合用于池，
		因为在这种情况下开销不能很好地摊销。让这些对象实现它们自己的空闲列表更有效。

		首次使用后不得复制池。
	*/
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	// 使用time.Now() 替换进入一个真正的日志读写器
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	b.Write(b.Bytes())
	bufPool.Put(b)

}

//// 自旋锁
type spinLock uint32

func (selflock *spinLock) Lock() {
	for !atomic.CompareAndSwapUnit32((*uint32)(selflock), 0, 1) {
		runtime.Gosched()
	}
}

func (selflock *spinLock) Unlock() {
	atomic.StoreUnit32((*uint32)(selflock), 0)
}

func NewSpinLock() sync.Locker {
	var lock spinLock
	return &lock
}
func main() {
	// OnceDoExample(10) // 循环10次，但 once.Do 只执行一次

	LogPoolDoExample(os.Stdout, "path", "/search?q=flowers")
}
