package main

import (
	"GolangTraining/02_package/stringutil" //相对导入 本地包时 需要 先 进行本地包的初始化 go mod init
	"fmt"
)

func main() {
	fmt.Println(stringutil.Reverse("!oG ,olleH"))
	fmt.Println(stringutil.MyName)
	fmt.Println(winniepooh.BearName)
}
