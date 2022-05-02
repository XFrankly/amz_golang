package main

import (
	"fmt"
	"sort"
)

func main() {
	s := []string{"Zeno", "John", "Al", "Jenny"}
	fmt.Println(s)
	//	sort.StringSlice(s).Sort()
	// String 分片
	sort.Sort(sort.StringSlice(s))
	fmt.Println(s)

}

// https://golang.org/pkg/sort/#Sort
