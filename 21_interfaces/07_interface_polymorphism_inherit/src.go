package main

import (
	"fmt"
)

type IChildHandles interface {
	AddAges()
	HandleAfterInits(handler IChildHandles)
}

type Bases struct {
	name string
	age  int
}

func (b *Bases) AddAges() {
	fmt.Println("base add Age .....")
}

func (b *Bases) reNames() {
	fmt.Println("rename ing...")
}

func (b *Bases) HandleAfterInits(handler IChildHandles) {
	fmt.Printf("before save base obj name:%s,age:%d\r\n", b.name, b.age)
	b.name = b.name + "aaaaa"
	handler.AddAges()
	fmt.Printf("real save base obj name:%s,age:%d\r\n", b.name, b.age)
}

type Child2 struct {
	*Bases
	price int
}

func (c *Child2) AddAges() {
	fmt.Println("child add age ...")
	c.age += 2
}

func main() {
	b := &Bases{
		name: "base xx",
		age:  100,
	}
	c := &Child2{
		Bases: b,
		price: 333,
	}
	b.HandleAfterInits(c)
	c.AddAges()
	c.reNames()
	fmt.Println(c.age)
}
