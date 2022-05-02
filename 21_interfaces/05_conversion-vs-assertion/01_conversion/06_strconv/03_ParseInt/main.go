package main

import (
	"fmt"
	"strconv"
)

func main() {

	//	ParseBool，ParseFloat，ParseInt和ParseUint将字符串转换为值：
	b, _ := strconv.ParseBool("true")
	f, _ := strconv.ParseFloat("3.1415", 64)
	i, _ := strconv.ParseInt("-42", 10, 64)
	ia,_ := strconv.Atoi("-11")
	u, _ := strconv.ParseUint("42", 10, 64)

	fmt.Println(b, f, i, ia, u)

	//	FormatBool, FormatFloat, FormatInt, and FormatUint 将值转为字符串:
	w := strconv.FormatBool(true)
	x := strconv.FormatFloat(3.1415, 'E', -1, 64)
	y := strconv.FormatInt(-42, 16)
	z := strconv.FormatUint(42, 16)

	fmt.Println(w, x, y, z)
}
