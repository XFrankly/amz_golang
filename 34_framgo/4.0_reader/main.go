package main

import (
	"bufio"
	"log"
	"os"
	"sync"
)

var (
	Logger    = log.New(os.Stderr, "INFO -", 18)
	code_path string
)

type MyChans struct {
	// *MyChan //
	Read <-chan map[int]any //interface{} // 只读通道  为 channel 通道创建一个 按索引查看的方法
	//all   chan map[int]interface{}   // 可读可写
	Input chan<- map[int]any //interface{} // 只写通道
	// maxsize int

}
type Queues struct {
	MaxSize int
	MC      *MyChans     //chan map[int]interface{} //  维护一个索引, 存储代码数据
	RWmu    sync.RWMutex // 读写锁
	mutex   sync.Mutex   // 互斥锁
	TokenC  *MyChans
}

func Readline(cpath ...string) *Queues {
	// cpath 可选代码路径，如果没有指定，将加载默认 的代码路径
	//按行 读取 代码文件，并存储在 channel 队列中
	// n 表示最多读 并存入chan 1000行
	if cpath != nil {
		code_path = cpath[0]
	}
	src, err := os.Open(code_path)
	if err != nil {
		log.Printf("error opening source file: %v", err)
	}

	var scanner1 *bufio.Scanner
	// var files *os.File
	scanner1 = bufio.NewScanner(src)
	scanner1.Split(bufio.ScanLines)
	var text []string
	Logg.Println("len before read text", len(text))
	for scanner1.Scan() {
		text = append(text, scanner1.Text()) //string(scanner1.Bytes()) //#
		// Logg.Println("text line: ", text, "len:", len(text))
	}

	defer src.Close()

	for x, each_in := range text {
		Logg.Println("No. ", x, "each in:", each_in, "len:", len(text))

	}

	Myqueues = PutQueue(Myqueues, text)
	return Myqueues
}
