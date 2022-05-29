package main

import (
	"fmt"
	"sync"
)

type NMyQueue struct {
	MaxSize int
	Putters []string //`form:"colors[]"`
}

var (
	Ninsqueue *NMyQueue
	Nmutex    sync.Mutex
	Nwg       = sync.WaitGroup{}
)

func NNewQueues(maxsize int) *NMyQueue {
	myq := &NMyQueue{
		MaxSize: maxsize,
		Putters: []string{},
	}
	return myq
}
func init() {

	Ninsqueue := NNewQueues(10)
	fmt.Sprintln(Ninsqueue.MaxSize) //#(im.Maxsize)
	Ninsqueue.NInitQueue()
}

func (self *NMyQueue) NPutQueue(data string) *NMyQueue {

	if len(self.Putters) < self.MaxSize {
		self.Putters = append(self.Putters, data)
	}
	// if C.length < self.MaxSize {
	// }

	return self
}

func (self *NMyQueue) NClearQueue() *NMyQueue {
	for n := range self.Putters {
		fmt.Sprintln(n)
		self.Putters = append(self.Putters)
	}
	//Nwg.Done()
	return self
}

func (self *NMyQueue) NInitQueue() *NMyQueue {
	// 删除ind 位置的元素
	for j := 0; j <= self.MaxSize; j++ {
		strs := fmt.Sprintf("a-%d", j)
		self.NPutQueue(strs)
		fmt.Println(self.Putters)
	}
	return self
}

func (self *NMyQueue) NCutQueue(ind int) *NMyQueue {
	// 删除ind 位置的元素
	if len(self.Putters) >= ind+1 {
		fmt.Println("del item at", ind)
		// self.Putters = append(self.Putters[:ind], self.Putters[ind+1:]...)

	}
	//Nwg.Done()

	fmt.Println(self.Putters[1:3])
	return self
}
func (self *NMyQueue) NCutAllQueue(ind int) *NMyQueue {
	// 删除ind 位置的元素
	if len(self.Putters) >= ind+1 {
		fmt.Println("del item at", ind)
		// self.Putters = self.Putters[ind:]

	}

	fmt.Println("NCutAllQueue:", self.Putters)
	//Nwg.Done()
	return self
}
func main() {
	Nwg.Add(3)
	Ninsqueue.NCutQueue(2)    // 删除第二个元素
	Ninsqueue.NCutAllQueue(2) //切除前2个元素

	Ninsqueue.NClearQueue()
	Nwg.Wait()
	fmt.Println(Ninsqueue.Putters)
	fmt.Println(Ninsqueue.Putters[:2])
}
