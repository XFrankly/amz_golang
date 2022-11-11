package bstdata

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
)

type node struct {
	numb *TreeNode
	prev *node
	next *node
}

type dlist struct {
	size int
	head *node
	tail *node
}

func makeDlist() *dlist {
	return &dlist{}
}

func (the *dlist) lesseq(n *node) (int, *node) {
	/// 排序 大于等于目标的结点
	if the.size <= 0 || n == nil || n.numb == nil {
		return 0, nil
	}
	currentNode := the.head
	for i := 0; i < the.size; i++ {
		if currentNode.numb.Key >= n.numb.Key {
			return i, currentNode
		} else {
			currentNode = currentNode.next
		}
	}
	/// 没有找到 比n 大的 返回尾节点位置
	return the.size - 1, nil
}

func (the *dlist) isNodeIn(n *node) (int, *node) {
	/// 节点是否存在于链表 大于等于目标的结点
	if the.size <= 0 || n == nil || n.numb == nil {
		return 0, nil
	}
	currentNode := the.head
	for i := 0; i < the.size; i++ {
		Logg.Printf("compare:%v with :%v\n", currentNode.numb, n.numb)
		if currentNode.numb.Key == n.numb.Key {
			return i, currentNode
		} else {
			currentNode = currentNode.next
			if currentNode == nil {
				break
			}
		}
	}
	/// 没有找到 比n 大的 返回尾节点位置
	return 0, nil
}

// 判断是否空链表
func (the *dlist) newNodeList(n *node) bool {

	if the.size == 0 {
		the.head = n
		the.tail = n
		n.prev = nil
		n.next = nil
		the.size += 1
		return true
	} else {
		Logg.Panic("not empty node list.")
	}
	return false
}

// 头部添加 节点
func (the *dlist) pushHead(n *node) bool {

	if the.size == 0 {
		return the.newNodeList(n)
	} else {
		the.head.prev = n
		n.prev = nil
		n.next = the.head
		the.head = n
		the.size += 1
		return true
	}

}

// 添加尾部节点
func (the *dlist) append(n *node) bool {

	if the.size == 0 {
		return the.newNodeList(n)
	} else {
		the.tail.next = n
		n.prev = the.tail
		n.next = nil
		the.tail = n
		the.size += 1
		return true
	}
}

// 修改链表 删除链表节点
func (the *dlist) delete(n *node) (bool, *node) {
	/// 节点是否存在于链表 大于等于目标的结点
	if the.size <= 0 || n == nil {
		return false, nil
	}
	currentNode := the.head
	for i := 0; i < the.size; i++ {
		if currentNode.numb.Key == n.numb.Key {
			//首节点删除
			currentNode.next.prev = nil
			//尾节点删除
			currentNode.prev.next = nil
			//中间节点删除
			oldNode := currentNode.prev
			newNode := currentNode.next
			currentNode.prev.next = newNode
			currentNode.next.prev = oldNode
			the.size -= 1
			return true, currentNode
		} else {
			currentNode = currentNode.next
			continue
		}
	}
	return false, nil
}

// 有序插入
func (the *dlist) pushback(n *node) bool {

	if n == nil || n.numb == nil {
		return false
	}
	currentNode := the.head
	if currentNode == nil {
		Logg.Printf("new list:%#v this :%+v\n", n, the)
		return the.newNodeList(n)
	} else {
		inDex, insertNode := the.lesseq(n)

		Logg.Printf("insert:%#v at :%+v size:%v\n", n.numb.Key, inDex, the.size)
		if inDex == 0 {
			return the.pushHead(n)
		} else if inDex == (the.size-1) && insertNode == nil {
			return the.append(n)
		}

		n.next = insertNode
		n.prev = insertNode.prev
		if insertNode.prev != nil {
			insertNode.prev.next = n
		}

		insertNode.prev = n
		the.size += 1
		return true
	}
}

func (the *dlist) display() []*TreeNode {
	/// 显示链表的值
	numbs := []*TreeNode{}
	node := the.head
	t := 0
	// Logg.Println(node.number)
	for node != nil {

		Logg.Println(node.numb)
		numbs = append(numbs, node.numb)
		t += 1
		fmt.Println("show the numb:", t)
		node = node.next

	}

	fmt.Println("length:", the.size)
	return numbs
}
