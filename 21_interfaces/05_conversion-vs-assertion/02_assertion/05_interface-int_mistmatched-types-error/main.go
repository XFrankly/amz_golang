package main

import "fmt"

func main() {
	var val interface{} = 7

	fmt.Println(val.(int) + 6) // 正确
	//fmt.Println(val + 6)  //尝试将接口与 整数直接相加
}
