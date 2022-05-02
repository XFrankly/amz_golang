package main

import (
	"fmt"
	"math"
)

type square struct {
	// 正方形结构体 边长
	side float64
}

// another shape
type circle struct {
	// 圆的结构体 直径
	radius float64
}

type shape interface {
	// 调用面积 函数
	area() float64
}

func (s square) area() float64 {
	// 实现计算正方形面积
	return s.side * s.side
}


func (c circle) area() float64 {
	// 实现 计算圆的 面积
	return math.Pi * c.radius * c.radius
}

func info(z shape) {
	// 调用 shape 接口
	fmt.Println("z shape", z)
	fmt.Println("z shape area", z.area())  // 正方形 与 圆  都有计算 面积的函数，这些函数都使用接口 area，
}

func main() {
	s1 := square{10}
	c := circle{5}
	info(s1)
	info(c)
}
