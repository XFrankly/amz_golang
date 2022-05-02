package main

import "fmt"

func zero(z int) {
	z = 0
}

func main() {
	x := 5
	zero(x)  // 参数 传递到 其他函数中 不会改变 参数的值
	fmt.Println(x) // x is still 5
}
