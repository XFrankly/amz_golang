package main

import "fmt"

func main() {
	str := []byte("hello")
	print(str)  // 0xc000070090
	fmt.Println(str)  //[104 101 108 108 111]
	// conversion: string to []bytes
}
