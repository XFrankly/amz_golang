package main

import "fmt"

func main() {
	var x rune = 'a' //符号 是int32的别名；通常在此语句中省略
	var y int32 = 'b'
	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(string(x))
	fmt.Println(string(y))

	// 转换 符号类型 ->  字符
	as := 'a'; print(as);   print(string(as)) //97 a
}
