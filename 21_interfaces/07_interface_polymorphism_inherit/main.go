package main

import (
	"fmt"
)

type IChildHandle interface {
	AddAge()
	HandleAfterInit(handler IChildHandle)
}

type Base struct {
	name string
	age  int
}

// func (b *Base) AddAge() {
// 	fmt.Println("base add age...")
// }

func (b *Base) reName() {
	fmt.Println("rename ...")
}

func (b *Base) HandleAfterInit(handler IChildHandle) {
	fmt.Printf("before save base obj name:%v, age:%v\n", b.name, b.age)
	b.name = b.name + "aaaa"
	handler.AddAge()
	fmt.Printf("real save base obj name:%v, age:%v\n", b.name, b.age)
}

type Child1 struct {
	*Base
	Price int
}

func (c *Child1) AddAge() {
	fmt.Println("child add age ...")
	c.age += 2
}

// func (c *Child1) NewFriend() {
// 	fmt.Println("add new friend...")
// 	c.age += 11
// 	c.Price += 100
// }
func main() {
	b := &Base{
		name: "base xx",
		age:  100,
	}
	c := &Child1{
		Base:  b,
		Price: 333,
	}
	b.HandleAfterInit(c)
	c.AddAge()
	c.reName()
	fmt.Println(c.age)
	// // b.NewFriend()
	// c.NewFriend()
	// fmt.Println(c.age)
}
