package main

import "fmt"

func main() {
	rs := []byte{'h', 'e', 'l', 'l', 'o'}
	fmt.Println(rs) //[104 101 108 108 111]
	fmt.Println(string(rs)) //hello
	// 转换: []bytes to string

}
