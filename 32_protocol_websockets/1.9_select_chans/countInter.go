package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

/*
统计 字符串中 元音字符出现次数
*/
const MaxCounts = 100

var (
	vowels    = []string{"a", "e", "i", "o", "u"}
	Muxs      sync.RWMutex
	SigCounts = 0
	Mc        MyCounts
	MyCounter = NewAdapterinters(MaxCounts)
)

type MyCounts int
type MyCountsruct struct {
	Read    <-chan int
	Write   chan<- int
	Total   int //存总值
	maxsize int
}

func NewMyCounts(size int) *MyCountsruct {
	c := make(chan int, size)
	return &MyCountsruct{
		Read:    c,
		Write:   c,
		Total:   0,
		maxsize: size,
	}
}

type CountInterface interface {
	putIns(i int) bool
	LookCount() int
	VowelsCounter(s string) bool
	BackCounter(t string, s string) int
}

// type CountAdapters struct {
// 	CountInterface
// }

func NewAdapterinters(size int) CountInterface {
	return NewMyCounts(size)
}

//PutIn 计算总数 并缓存最多最近50条数据
func (that *MyCountsruct) putIns(i int) bool {
	Muxs.Lock()
	defer Muxs.Unlock()
	if len(that.Write) >= 50 {
		that.Write = make(chan int, that.maxsize)
	}
	that.Write <- i
	that.Total += i

	return true
}

func (that *MyCountsruct) ReadOne() (int, bool) {
	Muxs.Lock()
	defer Muxs.Unlock()
	if len(that.Read) > 0 {

		return <-that.Read, true
	} else {
		return 0, false
	}
}

//LookCount 查询当前值
func (that *MyCountsruct) LookCount() int {
	return that.Total
}

func (that *MyCountsruct) number(s string) int {
	var t int
	for _, e := range vowels {
		t += strings.Count(s, e)
	}
	return t
}

func (that *MyCountsruct) VowelsCounter(s string) bool {
	return that.putIns(that.number(s))
}

func (that *MyCountsruct) BackCounter(t string, s string) int {
	if !BackEndType[SoundsType(t)] {
		return 0
	}
	f, _ := NewCountFunc(SoundsType(t), s)

	return f.LookCount()
}

func (tc *MyCounts) BackCounter(t string, s string) int {
	return MyCounter.BackCounter(t, s)
}

//叫声类型
type SoundsType string

const (
	// bird backend
	MyBirds SoundsType = "mybirds"
	// dog backend
	MyDogs SoundsType = "mydogs"
	// cat backend
	MyCats SoundsType = "mycats"
	// car backend
	MyCars SoundsType = "mycars"
	// tair backend
	MyTrain SoundsType = "mytrain"
)

//统计元音 输入函数
type Initialize func(c CountInterface, addrs string) (CountInterface, error)

var (
	ErrBackendNotSupported = errors.New("Backend storage not supported yet, please choose one of")
	BackEndType            = map[SoundsType]bool{
		MyBirds: true,
		MyDogs:  true,
		MyCats:  true,
		MyCars:  true,
		MyTrain: true,
	}

	initializers = map[SoundsType]Initialize{
		MyBirds: BirdSounds,
		MyDogs:  DefaultFunc,
		MyCats:  DefaultFunc,
		MyCars:  DefaultFunc,
		MyTrain: DefaultFunc,
	}

	//提示信息
	supportedBackend = func() string {
		keys := make([]string, 0, len(initializers))
		for k := range initializers {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)
		return strings.Join(keys, ", ")
	}()
)

func BirdSounds(c CountInterface, s string) (CountInterface, error) {
	if s == "" {
		s = "gua"
	}
	if c == nil {
		c = MyCounter
	}
	c.VowelsCounter(s)
	return c, nil
}

func DefaultFunc(counter CountInterface, s string) (CountInterface, error) {
	val := []reflect.Value{reflect.ValueOf(s)}
	reflect.ValueOf(counter).MethodByName("VowelsCounter").Call(val)
	return MyCounter, nil
}

//根据对象 类别 返回计数器
func NewCountFunc(backend SoundsType, addrs string) (CountInterface, error) {
	if init, exists := initializers[backend]; exists {
		// 调用 并更改计数值
		return init(MyCounter, addrs)
	}
	return nil, fmt.Errorf("%s %s", ErrBackendNotSupported.Error(), supportedBackend)
}

func NewCountFuncWithStruct(c CountInterface, backend SoundsType, addrs string) (CountInterface, error) {
	if c == nil {
		return NewCountFunc(backend, addrs)
	}
	if init, exists := initializers[backend]; exists {
		return init(c, addrs)
	}
	return nil, fmt.Errorf("%s %s", ErrBackendNotSupported.Error(), supportedBackend)
}

func main() {
	nf, err := NewCountFunc(MyBirds, "gua")
	fmt.Println(supportedBackend)
	fmt.Println("MyBirds", MyBirds, "create result:", nf, "err:", err)
	fmt.Println(nf.LookCount())
	fmt.Println(nf.BackCounter("mytrain", "bio"))
	fmt.Println(Mc.BackCounter("mytrain", "bio"))
}
