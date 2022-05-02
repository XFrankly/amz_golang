package main

import "fmt"

func main() {
	var x = 12
	var y = 12.1230123
	fmt.Println(y + float64(x))
	// 数据类型转换 int 转 float

	s, s1 := 11,  11.11111
	fmt.Println(s1 + float64(s))
}
