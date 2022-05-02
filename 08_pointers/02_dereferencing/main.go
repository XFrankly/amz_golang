package main

import "fmt"

func main() {

	a := 43

	fmt.Println(a)   // 43
	fmt.Println(&a)  // 0xc0000180a8
	fmt.Println(*&a) // 43

	var b = &a
	fmt.Println(b)   // 0xc0000180a8
	fmt.Println(*b)  // 43
	fmt.Println(&b)  // 0xc000006030
	fmt.Println(*&b) // 0xc0000180a8

	// b is an int pointer;
	// b points to the memory address where an int is stored
	// to see the value in that memory address, add a * in front of b
	// this is known as dereferencing
	// the * is an operator in this case
}
