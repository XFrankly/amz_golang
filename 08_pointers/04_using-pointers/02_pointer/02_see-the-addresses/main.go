package main

import "fmt"

func zero(z *int) {
	fmt.Println("from zero ", z)
	*z = 0  //指针加指针 指向了原值
}

func main() {
	x := 5
	fmt.Println(&x)
	zero(&x)
	fmt.Println(x) // x is 0
	fmt.Println(1+3333)
}
