package main

import "fmt"

func hello() {
	fmt.Print("hello ")
}

func world() {
	fmt.Println("world")
}

func forderfer() {
	for i := 0; i < 3; i++ {
		defer fmt.Print(i)
	}
}
func main() {
	defer world()
	hello()
	//逆序输出
	defer forderfer()
}
