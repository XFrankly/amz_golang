package main

import (
	"fmt"
	"sort"
)

func main() {
	// 字符串排序
	s := []string{"Zeno", "John", "Al", "Jenny"}
	sort.Strings(s)
	fmt.Println(s)
}
