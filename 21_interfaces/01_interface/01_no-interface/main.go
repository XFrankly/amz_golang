package main

import "fmt"

type square struct {
	side float64
}

func (z square) area() float64 {
	return z.side * z.side
}

func main() {
	// 求平方乘机，即面积
	s := square{10}
	fmt.Println("Area: ", s.area())
}
