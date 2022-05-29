package main

/*
https://groups.google.com/g/golang-nuts/c/hLpwaPOPV2I
这个线程为通道提供了一些包装器，以启用 Peek() 函数，但这更像是一种解决方法。

type PeekChanInt struct {
        in <-chan int
        out chan int
}

	使用 channel 实现一个队列，
		可以 添加元素到队列
		可以 按索引下标 从队列获取某个元素，并从队列删除
		可以 按 索引下标 访问队列的 元素
		可以 查看索引，这些索引标识队列 剩下那些元素

		当 访问队列中的元素时，更新该元素到 队尾

*/
import (
	"fmt"
	"sync"
	"time"
)

type MyChans struct {
	// *MyChan //
	read <-chan map[int]interface{} // 只读通道  为 channel 通道创建一个 按索引查看的方法
	//all   chan map[int]interface{}   // 可读可写
	input chan<- map[int]interface{} // 只写通道
	// maxsize int

}

// MyChan的构造函数
func FuncMyChans(maxsize int) *MyChans {
	// 只读 只写 分开
	var MyC = make(chan map[int]interface{}, maxsize)
	ret_mychan := &MyChans{
		read:  MyC,
		input: MyC,
		// maxsize: maxsize,
	}
	return ret_mychan
}

type Queues struct {
	MaxSize int
	MC      *MyChans     //chan map[int]interface{} //  维护一个索引
	mu      sync.RWMutex // 读写锁
	mutex   sync.Mutex   // 互斥锁
}

func NewMyQueues(maxsize int) *Queues {
	//  构造Queues
	nc := FuncMyChans(maxsize)
	myq := &Queues{
		MaxSize: maxsize,
		MC:      nc,
	}
	return myq
}

var (
	squeue *Queues
	mutx   sync.Mutex
	// 用于 存储 原始数据
	// C   chan []interface{} //map[int]interface{}
	wg1  = sync.WaitGroup{}
	wg12 = sync.WaitGroup{}
)

func (self *Queues) Put(data map[int]interface{}) {
	// 向管道添加数据
	self.mu.RLocker().Lock()
	fmt.Println("try lock ", self.mu.RLocker())
	defer self.mu.RLocker().Unlock()
	// self.mu.Lock()
	// defer self.mu.Unlock()
	if len(self.MC.input) < self.MaxSize {
		self.MC.input <- data
	} else {
		fmt.Println("There are no place to append more items to Chan, Please delete first.")
	}
}

func (self *Queues) ReadChannel() map[int]interface{} {
	// 从管道取数据
	if len(self.MC.read) > 0 {
		// self.Indexs = self.Indexs[1:] // 切除第一个元素
		// tl := self.mu.TryLock()
		self.mu.RLocker().Lock()
		fmt.Println("try lock ", self.mu.RLocker())
		// defer self.mu.Unlock()
		defer self.mu.RLocker().Unlock()
		return <-self.MC.read
	} else {
		panic("There are no data in Chan to read, Please put first.")
	}
}

func (self *Queues) Get(ind int) map[int]interface{} {
	// 按索引 取某个值，环形队列, 并从 channel删除
	var ret_value map[int]interface{}
	//fmt.Println("chan self.MC work normal len MC.read", len(self.MC.read), "ret_value: ", ret_value)
	for i := 0; i <= len(self.MC.read); i++ {
		x := self.ReadChannel()
		for k, _ := range x {
			if k == ind {
				// fmt.Println("Origin index:", k, "value: ", v, "equal ind:", k == ind, "\n")
				ret_value = x
			} else {
				// fmt.Println("Origin:", k, "value: ", v, "equal ind:", k == ind, "\n")
				self.Put(x) //重新添加到 channel
				// fmt.Println("after of top of buffer:", x, "and len of MC", len(self.MC.input))

			}
		}
	}
	return ret_value
}

func (self *Queues) Delete(ind int) {
	/// 删除 通道中的 目标元素 map[int]string 中 int相同的
	for i := 0; i <= len(self.MC.read); i++ {
		//fmt.Println("CircleMyChannel Times Number:", i, "read length:", len(self.MC.read), "input length:", len(self.MC.input))
		n := self.ReadChannel() //<-self.MC.read
		for k, v := range n {
			//fmt.Println(" index:", k, "value: ", v, "equal ind:", "\n")
			if k != ind { // 删除 目标 元素，非目标元素，则保留
				// self.MC.input <- n
				self.Put(n)
			} else {
				fmt.Println("delete success item:", ind, "value:", v)
			}
		}

	}
	wg12.Done()
}

func (self *Queues) Indexs() []int {
	///  查询 channel 元素索引
	var index_items []int //, self.MaxSize)
	fmt.Println("length channel", len(self.MC.read))
	for j := 0; j < len(self.MC.read); j++ {
		// 获取
		datas := self.ReadChannel()
		for k, v := range datas {
			if v != nil {
				// 检测 格式 并添加到 索引 队列
				index_items = append(index_items, k)
			} else {
				panic("value not exist.")
			}
		}
		// 放回
		self.Put(datas)
	}
	return index_items
}

func (self *Queues) AccessItem(ind int) {
	/// 查找 通道中的 目标元素 map[int]string 中 int相同的， 并更新它到 channel 顶部
	var dist_item map[int]interface{}
	// fmt.Println("dist item", dist_item)
	n := self.Get(ind)
	for k, _ := range n { // 校验
		// fmt.Println(" index:", k, "value: ", v, "equal ind:", "\n")
		if k != ind { // 删除 目标 元素，非目标元素，则保留
			panic("get the wrong item.")
		} else {
			fmt.Println("update success item:", ind)
			dist_item = n
		}
	}

	fmt.Println("dist item>?", dist_item)
	self.Put(dist_item)
	wg12.Done()
}

func (self *Queues) InitQueues() {
	var inputInterArray = make([]interface{}, self.MaxSize)
	// inputInterArray = []interface{}{}
	for j := 0; j < self.MaxSize; j++ {
		// self.PutChannel()
		indexMap := make(map[int]interface{}) //, self.MaxSize)
		// var indexMap map[int]interface{}
		datas := fmt.Sprintf("c-%d", j)
		// 原始数据 索引和 数据 以方便 与 chan的数据 比对
		inputInterArray[j] = datas
		indexMap[j] = datas
		// 填充到 实例的 chan
		self.Put(indexMap)

	}
	wg12.Done()
}

func main() {
	// imq := NewMyQueues(10)
	// wg12.Add(2)
	// imq.InitQueues()
	// imq.AccessItem(2)
	// wg12.Wait()

	// fmt.Println("imq info:", imq.Indexs(), imq.MaxSize, len(imq.MC.input))
	// // fmt.Println("read channel item 2:", imq.Get(2), imq.Indexs())
	// // wg12.Add(1)
	// // imq.AccessItem(2)
	// // wg12.Wait()
	// fmt.Println("Indexs channel item ", imq.Indexs())

	// wg12.Add(1)
	// imq.Delete(3)
	// wg12.Wait()
	// fmt.Println("imq info:", imq.Indexs(), imq.MaxSize, len(imq.MC.input))
	// close(imq.MC.input)

	imq2 := NewMyQueues(5)
	for i := 0; i <= imq2.MaxSize; i++ {
		indexMap := make(map[int]interface{}, imq2.MaxSize)
		indexMap[i] = fmt.Sprintf("a-%d", i)
		imq2.Put(indexMap)
		// indexMap[i+1] = fmt.Sprintf("b-%d", i)
		// imq2.Put(indexMap)
	}

	// indexMap[3] = "b"
	// imq2.Put(indexMap)   // 同时 put 多键到单键map将失败，因为map 应该是单键map
	// imq2.Put(indexMap)
	// imq2.Put(indexMap)
	// imq2.Put(indexMap)
	// imq2.Put(indexMap)
	// imq2.Put(indexMap)
	fmt.Println("index map", imq2.Indexs(), len(imq2.MC.read))
	d3 := imq2.Get(3)
	for k, v := range d3 {
		fmt.Println("d3 k:", k, "v:", v)
	}
	for {
		if len(imq2.MC.read) > 0 {
			dd := imq2.ReadChannel()
			fmt.Println("read success:", dd, len(imq2.MC.read))
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

}
