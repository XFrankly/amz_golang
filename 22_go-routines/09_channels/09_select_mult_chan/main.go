package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func selectChan() {
	var c1, c2, c3 chan int
	var i1, i2 int
	select {
	case i1 = <-c1:
		fmt.Printf("received", i1, "from c1 \n")
	case c2 <- i2:
		fmt.Printf("sent:", i2, "to c2\n")
	case i3, ok := (<-c3): // same as: i3, ok:<-c3
		if ok {
			fmt.Printf("receive:", i3, " from c3\n")

		} else {
			fmt.Printf("c3 is closed \n")
		}
	default:
		//// 一定会执行这个 default，不会去等待 c1 c2
		fmt.Printf("no communication\n")

	}
}

/////// select的随机执行, 只要 chan 没有关闭就会一直执行
func Chann(ch chan int, stopCh chan bool) {
	for j := 0; j < 10; j++ {
		ch <- j
		time.Sleep(1 * time.Second)
	}
	stopCh <- true
}

func RandSelect() {
	ch := make(chan int)
	c := 0
	stopCh := make(chan bool)

	go Chann(ch, stopCh)

	for {
		select { /// 同一个 ch 随机选择一个
		case c = <-ch:
			fmt.Println("receive C", c)
		case s := <-ch:
			fmt.Println("Receive s", s)
		case sp := <-stopCh:
			fmt.Println("stopch goto end:", sp)
			goto end
		}
	}
end: //// 相当与 python的 while.... else， 在选择外部执行
	fmt.Println("its really end....")
}

////////////对多个通道进行选择
/// 每个通道将在一段时间收到一个值，以模拟例如在并发 goroutine中执行的阻塞RPC操作
//// 我们将使用select 同时等待这两个值，每个值到达时打印它们，总执行时间 2秒
//// 因为 1秒 和 2秒 sleep 同时执行
var (
	ge errgroup.Group
)

func TwoChanSelect() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		// ge.Go(func() error {
		time.Sleep(1 * time.Second)
		c1 <- "one"
		// return nil
		// })
	}()

	go func() {
		// ge.Go(func() error {
		time.Sleep(2 * time.Second)
		c2 <- "two"
		// return nil
		// })
	}()

	go func() {
		// ge.Go(func() error {
		time.Sleep(1 * time.Second)
		c2 <- "three"
		// return nil
		// })
	}()

	ge.Wait()

	for i := 0; i < 3; i++ {
		select { ////
		case msg1 := <-c1:
			fmt.Println("received:", msg1)
		case msg2 := <-c2:
			fmt.Println("received:", msg2)
			// default:
			// 	fmt.Println("Nothing at here.")
		}
	}
}

func main() {
	// selectChan()
	RandSelect()
	TwoChanSelect()
}
