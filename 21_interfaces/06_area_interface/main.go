package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
	attr() float64
}

type shape2 interface {
	perim() float64
}

type circle struct {
	radius float64
}

type rect struct {
	width  float64
	height float64
}

func (r *rect) area() float64 {
	return r.width * r.height
}
func (r *rect) attr() float64 {
	return r.height
}

func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c *circle) attr() float64 {
	return c.radius
}
func getArea(s shape) float64 {
	return s.area()
}

func main() {
	c1 := circle{4.5}           /// 任务对象1
	r1 := rect{5, 7}            /// 任务对象 2
	shapes := []shape{&c1, &r1} /// 创建任务队列
	for x, shape := range shapes {
		fmt.Printf("get info:x %+v\n shape %+v\n", x, shape) //// 获得的数据信息
		fmt.Printf("obj: %+v area:%+v\n", shape.attr(), getArea(shape))
	}

}
