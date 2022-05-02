package main

import "fmt"

type square struct {
	//结构体 正方形 边长
	side float64
}

func (z square) area() float64 {
	// 实例化结构体 ，求其面积
	return z.side * z.side
}

type shape interface {
	// 面积的 接口
	area() float64
}

func info(z shape) {
	// 调用接口
	fmt.Println(z)
	fmt.Println(z.area())
}

func main() {
	//
	s := square{10}
	fmt.Printf("%T\n", s)
	info(s)
}
