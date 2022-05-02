package main

import "fmt"

func main() {
	mySlice := make([]int, 1)
	fmt.Println(mySlice[0])
	mySlice[0] = 7
	fmt.Println(mySlice[0])
	mySlice[0]++    //直接对切片元素 做运算
	fmt.Println(mySlice[0])
	fmt.Println(mySlice)

}
