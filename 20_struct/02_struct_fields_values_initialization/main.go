package main

import "fmt"

type person struct {
	//声明自己的类型。
	// 通过将一组固定的唯一字段组合在一起来声明结构类型。
	// 结构中的每个字段都以已知类型声明。
	// 这可以是内置类型，也可以是其他用户定义类型。
	// 声明类型后，就可以从该类型创建值
	// 当我们声明变量时，变量代表的值总是被初始化。
	// 可以使用特定值初始化该值，也可以将其初始化为零值
	// 对于数字类型，零值为0；对于数字类型，零值为0。对于字符串，它将为空；
	// 对于布尔值，这将是错误的。
	// 在结构的情况下，零值将应用于该结构中的所有不同字段。
	// 每当创建变量并将其初始化为零值时，使用关键字var都是惯用的。
	// 保留使用关键字var作为指示将变量设置为零值的方式。
	// 如果将变量初始化为零值以外的其他值，
	// 然后将短变量声明运算符与struct文字一起使用
	first string
	last  string
	age   int
}

func main() {
	p1 := person{"James", "Bond", 20}
	p2 := person{"Miss", "Moneypenny", 18}
	fmt.Println(p1.first, p1.last, p1.age)
	fmt.Println(p2.first, p2.last, p2.age)
}
