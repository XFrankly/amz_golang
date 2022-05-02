package main

import "fmt"

type person struct {
	name string
	age  int
}

type SinglyListNode struct {
	// 单节点 链表
	val  int
	next *SinglyListNode
}

func main() {
	p1 := SinglyListNode{201} //直接 指向结构体地址 使用
	fmt.Println(p1)
	fmt.Printf("%T\n", p1)
	fmt.Println(p1.val)
	// fmt.Println(p1.age)
}
