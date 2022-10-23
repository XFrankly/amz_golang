package main

import (
	"fmt"
	"time"
)

var x = 10

var a = []int{1, 2, 4, 5, 7, 8}
var a2 = []int{3, 5, 6, 9, 10, 11}
var a3 = []int{} //make([]int, 10, 100)

func main() {
	fmt.Println("package var x", x)
	start := time.Now() //time.NewTimer()
	aa2 := a2[:]
	for i, v := range a {
		for j, va := range aa2 {
			fmt.Println(j, va, "in", a2)
			if v < va {
				a3 = append(a3, v)
				break
			} else {
				a3 = append(a3, aa2[0]) // 添加aa2的第一个元素
				aa2 = aa2[1:]           // 删除第一个元素
			}
		}
		if len(aa2) == 0 {
			a3 = append(a3, a[i:]...)
		}
	}
	if len(aa2) != 0 {
		a3 = append(a3, aa2...)
	}
	fmt.Println(a2, a3)
	end := time.Now()
	fmt.Println("Fininal a3 is\n", a3, start, "\n", end)
	fmt.Println("cost time", end == start)
}
