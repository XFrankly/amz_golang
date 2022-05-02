package main

import "fmt"

func main() {
	var val interface{} = 7   //接口为 int
	fmt.Printf("%T\n", val)
}
