package main

import (
	"fmt"
	"math"
)

type circle struct {
	//圆的 结构体
	radius float64
}

type square struct {
	// 正方形结构体
	side float64
}

type shape interface {
	// 面积 接口
	area() float64
}

func (c circle) area() float64 {
	// 实现 计算圆的面积 的函数
	return math.Pi * c.radius * c.radius
}

func (s square) area() float64 {
	// 实现 计算正方形 面积 的函数
	return s.side * s.side
}

func info(z shape) {
	// 调用 计算面积 的函数
	fmt.Println(z)
	fmt.Println(z.area())
}

//一种采用INTERFACE TYPE形状的新方法
func totalArea(shapes ...shape) float64 {
	// 定义 面积 变量
	var area float64
	for _, s := range shapes {
		// 计算 总数
		area += s.area()
	}
	return area
}

func main() {
	s := square{10}
	c := circle{5}
	info(s)
	info(c)
	fmt.Println("Total Area: ", totalArea(c, s)) // 计算 总 面积
}
