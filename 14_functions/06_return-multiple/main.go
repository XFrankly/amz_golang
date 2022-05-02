package main

import "fmt"

func main() {
	fmt.Println(greet("Jane ", "Doe "))
}

func greet(fname, lname string) (string, string) {
	return fmt.Sprint(fname, lname), fmt.Sprint(lname, fname)  //使用默认格式为其操作数设置Sprint格式，并返回结果字符串
}
