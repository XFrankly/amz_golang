package main

import (
	"fmt"
	"strconv"
)

func main() {
	var x int = 5
	ix := strconv.Itoa(x)
	fmt.Printf("%T %T\n", x, ix)
	str := "Hello world " + ix // int to ascii
	fmt.Println(str)
}
