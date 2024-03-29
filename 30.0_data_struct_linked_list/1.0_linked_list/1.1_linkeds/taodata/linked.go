package taodata

/*
// 随机输入数字，形成链表
// 要求
// 1、 链表为双向链表
// 2、 链表为顺序排列的链表

// 例：
// 输入： 9，2，5，6，7，2，6
// 输出： 2<->2<->5<->6<->6<->7<->9

// 输入 10
// 输出 ：2<->2<->5<->6<->6<->7<->9<->10

// 输入 3
// 输出 2<->2<->3<->5<->6<->6<->7<->9<->10
*/
import (
	"fmt"
	"log"
	"os"
)

var (
	Logg = log.New(os.Stderr, "INFO - ", 18)
)

type node struct {
	number  int
	yaobian [][]int
	prev    *node
	next    *node
}

type dlist struct {
	lens int
	head *node
	tail *node
}

func makeDlist() *dlist {
	return &dlist{}
}

func (the *dlist) lesseq(n *node) (int, *node) {
	/// 排序 大于等于目标的结点
	if the.lens <= 0 || n == nil {
		return 0, nil
	}
	currentNode := the.head
	for i := 0; i < the.lens; i++ {
		if currentNode.number >= n.number {
			return i, currentNode
		} else {
			currentNode = currentNode.next
		}
	}
	/// 没有找到 比n 大的
	return the.lens - 1, nil
}

// 判断是否空链表
func (the *dlist) newNodeList(n *node) bool {

	if the.lens == 0 {
		the.head = n
		the.tail = n
		n.prev = nil
		n.next = nil
		the.lens += 1
		return true
	} else {
		Logg.Panic("not empty node list.")
	}
	return false
}

// 头部添加 节点
func (the *dlist) pushHead(n *node) bool {

	if the.lens == 0 {
		return the.newNodeList(n)
	} else {
		the.head.prev = n
		n.prev = nil
		n.next = the.head
		the.head = n
		the.lens += 1
		return true
	}

}

// 添加尾部节点
func (the *dlist) append(n *node) bool {

	if the.lens == 0 {
		return the.newNodeList(n)
	} else {
		the.tail.next = n
		n.prev = the.tail
		n.next = nil
		the.tail = n
		the.lens += 1
		return true
	}
}

// 有序插入
func (the *dlist) pushback(n *node) bool {

	if n == nil {
		return false
	}
	currentNode := the.head
	if currentNode == nil {

		return the.newNodeList(n)
	} else {
		inDex, insertNode := the.lesseq(n)
		if inDex == 0 {
			return the.pushHead(n)
		} else if inDex == (the.lens-1) && insertNode == nil {
			return the.append(n)
		}
		Logg.Printf("insert at :%+v\n", inDex)

		n.next = insertNode
		n.prev = insertNode.prev
		//// 很容易失败
		if insertNode.prev != nil {
			insertNode.prev.next = n
		}

		insertNode.prev = n
		the.lens += 1
		return true
	}
}

func (the *dlist) display() []int {
	/// 显示链表的值
	numbs := []int{}
	node := the.head
	t := 0
	// Logg.Println(node.number)
	for node != nil {

		Logg.Println(node.number, node.yaobian)
		numbs = append(numbs, node.number)
		t += 1
		if t >= the.lens {
			break
		}

		node = node.next
	}

	fmt.Println("length:", the.lens)
	return numbs
}

// func main() {
// 	dlist := makeDlist()
// 	slit := []int{9, 2, 5, 6, 7, 2, 6, 10, 3}
// 	for _, i := range slit {
// 		node := &node{number: i}
// 		// node.prev = node
// 		// node.next = node
// 		dlist.pushback(node)
// 	}

// 	dlist.display()
// 	Logg.Println()
// 	dlist.pushback(&node{number: 123})
// 	dlist.display()

// 	Logg.Println()
// 	dlist.pushback(&node{number: 323})
// 	dlist.display()
// 	Logg.Println()
// 	dlist.pushback(&node{number: 0})
// 	dlist.display()
// 	Logg.Println()
// 	dlist.pushback(&node{number: 1})
// 	dlist.display()
// 	Logg.Println()
// }
