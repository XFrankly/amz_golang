package main

import "fmt"

type foo int

func main() {
	var myAge foo //引用别名，不可这样操作
	myAge = 44
	fmt.Printf("%T %v \n", myAge, myAge)
}
