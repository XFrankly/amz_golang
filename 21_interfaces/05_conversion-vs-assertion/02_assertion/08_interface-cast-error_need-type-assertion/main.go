package main

import "fmt"

func main() {
	rem := 7.24
	fmt.Printf("%T\n", rem)
	fmt.Printf("%T\n", int(rem))

	var val interface{} = 7
	fmt.Printf("%T\n", val)
	fmt.Printf("%T\n", int(val))  // 转换接口的 int 类型数据 需要 先取出值
		fmt.Printf("%T\n", val.(int))
}
