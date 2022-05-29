package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type spinLock struct {
	owner int
	count int
}

func (spLock *spinLock) Lock() {
	me := GetGoroutineId()
	if spLock.owner == me { //如果当前线程获得锁，计数加1并返回
		spLock.count++
		return
	}
	for !atomic.CompareAndSwapUint32((*uint32)(spLock), 0, 1) {
		runtime.Gosched()
	}
}

func GetGoroutineId() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func NewSpinLock() sync.Locker {
	var lock spinLock
	return &lock
}
