package main

import "fmt"

func main() {
	name := "Sydney"
	str, ok := name.(string) //非法代码， 尝试将 字符串 当作接口用
	if ok {
		fmt.Printf("%q\n", str)
	} else {
		fmt.Printf("value is not a string\n")
	}
}
