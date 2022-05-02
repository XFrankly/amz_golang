package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 字符数字转换 int 或 float
	var x = "12"
	var y = 6
	z, _ := strconv.Atoi(x)
	fmt.Printf("%T %T \n", x, y)
	fmt.Println(y + z)

	s := "1.1314";fmt.Printf("%T \n", s);z1, _:= strconv.ParseFloat(s, 64)

	fmt.Println(float64(z) + z1)
}
