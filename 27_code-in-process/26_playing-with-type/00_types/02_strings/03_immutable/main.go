package main

import "fmt"

//不可变地址 将不同的 字符 赋予 相同的地址
func main() {
	intro := "Four score and seven years ago...."
	fmt.Println(intro)
	fmt.Println(&intro)
	intro = "Hahahaha!"
	fmt.Println(intro)
	fmt.Println(&intro)
	//  the below is invalid
	//	intro[0] = 70
	//	fmt.Println(intro)
	//	fmt.Println(&intro)
}
