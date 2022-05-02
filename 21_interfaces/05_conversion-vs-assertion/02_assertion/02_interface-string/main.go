package main

import "fmt"

func main() {
	var name interface{} = "Sydney"
	str, ok := name.(string)  // 正常调用接口
	if ok {
		fmt.Printf("%T\n", str)
	} else {
		fmt.Printf("value is not a string\n")
	}
}
