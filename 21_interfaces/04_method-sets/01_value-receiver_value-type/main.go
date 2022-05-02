package main

import (
	"fmt"
	"math"
)

type circle struct {
	//圆 结构体
	radius float64
}

type shape interface {
	//接口 不可以有打印
	area() float64
}

func (c circle) area() float64 {
	fmt.Println("circle")
	return math.Pi * c.radius * c.radius
}

func info(s shape) {
	fmt.Println("area", s.area())
}

func main() {
	c := circle{5}
	info(c)
}
