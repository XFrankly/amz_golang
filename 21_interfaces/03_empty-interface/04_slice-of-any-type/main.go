package main

import "fmt"

type animal struct {
	//结构体 动物
	sound string
}

type dog struct {
	//结构体 嵌套子类
	animal
	friendly bool
}

type cat struct {
	//结构体 嵌套子类
	animal
	annoying bool
}

func main() {
	fido := dog{animal{"woof"}, true}
	fifi := cat{animal{"meow"}, true}
	shadow := dog{animal{"woof"}, true}
	critters := []interface{}{fido, fifi, shadow}
	fmt.Println(critters)
}
