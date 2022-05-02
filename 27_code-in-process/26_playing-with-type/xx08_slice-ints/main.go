package main

import "fmt"

type myType []int //声明一个切片

func main() {
	var x myType = []int{32, 44, 57} //赋予切片一个值
	fmt.Println(x)
	fmt.Printf("%T\n", x)
	fmt.Printf("%T\n", []int(x))
}
