package main

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	mutex sync.Mutex // 互斥锁
)

type M struct {
	Map map[string]string
}

// 添加数据
func (m *M) Set(key, value string) {
	mutex.Lock()
	m.Map[key] = value
	mutex.Unlock()
}

// 获取数据
func (m *M) Get(key string) string {
	return m.Map[key]
}

func Do_thing() {
	c := M{Map: make(map[string]string)}
	wg := sync.WaitGroup{}
	for i := 0; i < 21; i++ {
		wg.Add(1)
		go func(n int) {
			k, v := strconv.Itoa(n), strconv.Itoa(n)
			c.Set(k, v)
			fmt.Printf("k:=%v, v:=%v\n", k, c.Get(k))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func main() {
	Do_thing()
}
