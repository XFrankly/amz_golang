package main

import (
	"fmt"
	"sort"
)

func main() {
	n := []int{5, 2, 6, 3, 1, 4}
	sort.Ints(n)
	fmt.Println(n)


	type people []string
	studyGroup := people{"Zeno", "John", "Al", "Jenny"}

	s := []string{"Zeno", "John", "Al", "Jenny"}

	sort.Sort(sort.Reverse(sort.StringSlice(studyGroup)))  // 字符 分片组 逆序
	sort.Sort(sort.Reverse(sort.StringSlice(s)))  // 字符分片组 逆序

	n1 := []int{7, 4, 8, 2, 9, 19, 12, 32, 3}
	sort.Sort(sort.Reverse(sort.IntSlice(n1)))  // int 分片组 逆序

	fmt.Println(n)
	fmt.Println(s)
	fmt.Println(studyGroup)

}
