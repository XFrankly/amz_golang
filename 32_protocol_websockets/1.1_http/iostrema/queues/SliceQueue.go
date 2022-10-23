package main

import (
	"fmt"
	"sync"
)

/*
	使用 slice 实现一个队列，
		可以 添加元素到队列
		可以 按索引下标 从队列获取某个元素，并从队列删除
		可以 按 索引下标 访问队列的 元素
		可以 查看索引，这些索引标识队列 剩下那些元素

		当 访问队列中的元素时，更新该元素到 队尾

*/

type MyQueue struct {
	MaxSize int
	Putters []string
}

func NewQueues(maxsize int) *MyQueue {

	myq := &MyQueue{
		MaxSize: maxsize,
		Putters: []string{},
	}
	return myq
}

var (
	insqueue *MyQueue
	mutex    sync.Mutex
	// 用于 存储 原始数据
	// C   chan []interface{} //map[int]interface{}
	wg  = sync.WaitGroup{}
	wg2 = sync.WaitGroup{}
)

///////////////////////// slice 实现
func (self *MyQueue) PutQueue(data string) *MyQueue {

	if len(self.Putters) < self.MaxSize {
		self.Putters = append(self.Putters, data)
	} else {
		fmt.Println("There are no place to append more items, delete at first.")
	}

	return self
}

func (self *MyQueue) ClearQueue() *MyQueue {
	self.Putters = []string{}
	return self
}

func (self *MyQueue) InitQueue() *MyQueue {
	for j := 0; j <= self.MaxSize; j++ {
		strs := fmt.Sprintf("a-%d", j)
		self.PutQueue(strs)
		fmt.Println(self.Putters)
		fmt.Println(len(self.Putters))
	}
	defer wg.Done()
	return self
}

func (self *MyQueue) GetItem(ind int) string {
	// 获取当前文字的值, 并合并 队列
	var item string
	if len(self.Putters) > ind {
		item = self.Putters[ind]
		self.Putters = append(self.Putters[:ind], self.Putters[ind+1:]...)
		// <-C // C 通道也取出一个值

	}
	return item
}

func (self *MyQueue) AccessItem(ind int) string {
	item := self.GetItem(ind)
	self.Putters = append(self.Putters, item)
	return item
}

func main() {
	C := make(chan string, 10)
	im := NewQueues(10)
	// defer close(im.MC.input)
	fmt.Sprintln(im.MaxSize) //#(im.Maxsize)
	fmt.Println("befer putter length:", len(im.Putters))
	wg.Add(1)
	// 初始化
	go im.InitQueue()
	wg.Wait() // wait 之后 关闭后 再读取 Channel通道，不然可能死锁deadlock

	// close(C)

	// 取值 并删除
	fmt.Println(im.GetItem(2))
	// 查询库存
	fmt.Println("after GetItem storage remain after get NO.2 Item:", im.Putters)
	// 访问值，并更新队列
	fmt.Println(im.AccessItem(2))
	// 查询库存
	fmt.Println("after AccessItem storage remain access item :", im.Putters)
	// 重置为空
	im.ClearQueue()
	// 查询库存
	fmt.Println("storage remain after clear storage:", im.Putters)

	close(C) // 写关闭

	// C2 := make(chan int, 2)
	// defer close(C2)
	// C2 <- 12
	// fmt.Println("len C2:", len(C2)) // len: 1
	// C2 <- 13
	// fmt.Println("len C2:", len(C2)) // len: 2
	// //C2 <- 4 //死锁 超过channel的缓冲区限制的 存入 fatal error: all goroutines are asleep - deadlock!
	// fmt.Printf("len C2:%d, get one:%d, after len get:%d \n", len(C2), <-C2, len(C2))
	// fmt.Printf("get one:%d, and length:%d \n", <-C2, len(C2))
	// //<-C2 //死锁 channel为空，却期望取得值 fatal error: all goroutines are asleep - deadlock!
	// fmt.Println("buf cap chan:", cap(C2))

	///////////// 只读chan
	// nm := NewMyChan(10)

	// indexM := make(map[int]interface{}, nm.maxsize)
	// indexM[1] = "a"
	// nm.input <- indexM
	// fmt.Println("nm", len(nm.input), len(nm.read), cap(nm.input), cap(nm.read))
	// fmt.Printf("read: %v \n", <-nm.read)
	// fmt.Println("len:", len(nm.input), len(nm.read))

	// wg2.Add(1)
	// go im.CircleChannel()
	// wg2.Wait()

	// wg2.Add(1)
	// go im.CircleMyChannel(0)
	// wg2.Wait()

	// wg2.Add(1)
	// go im.DelMyChannelInd(0)
	// wg2.Wait()
	// fmt.Println("im length channel:", len(im.MC.read), len(im.MC.input))
	// fmt.Println(im.MC.read)

	// wg2.Add(1)
	// go im.UpdateMyChannelInd(8)
	// wg2.Wait()

	// fmt.Println("length channel:", len(im.MC.read), im.Indexs)
	// close(im.MC.input)

}
