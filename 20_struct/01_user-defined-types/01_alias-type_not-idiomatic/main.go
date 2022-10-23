package main

import "fmt"

//类型定义
type foo int

//类型别名
type bar = int

func main() {
	var myAge foo //引用类型，不可这样操作
	var myBar bar
	myAge = 44
	myBar = 66

	fmt.Printf("%T %v \n", myAge, myAge) //main.foo 44
	fmt.Printf("%T %v \n", myBar, myBar) // int 66
}
