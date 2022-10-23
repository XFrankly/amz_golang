package main

import (
	"fmt"
	"sync"
)

type ipMap struct {
	sm   sync.Mutex
	ips  map[string]int
	stat struct {
		Size, Adds, Lookups, Dels int
	}
}

func NewIpMaps() *ipMap {
	retMap := &ipMap{}
	retMap.sm = sync.Mutex{}
	retMap.ips = make(map[string]int, 10)
	retMap.stat = struct {
		Size    int
		Adds    int
		Lookups int
		Dels    int
	}{Size: 10, Adds: 10, Lookups: 10, Dels: 10}
	return retMap
}

func main() {
	nm := NewIpMaps()
	fmt.Println(nm.ips)
	fmt.Println(nm.sm)
	fmt.Println(nm.stat.Size)
}
