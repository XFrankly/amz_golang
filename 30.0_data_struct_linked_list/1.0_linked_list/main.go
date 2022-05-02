package main

import "fmt"

type Node struct {
	Data int   // 包含数据的整数类型
	Next *Node // 保存下一个节点的内存地址 next 字段
}

type LinkedList struct {
	Length int   // 链表的长度，头节点，尾节点
	Head   *Node // 存放链表头或首节点的 内存地址
	Tail   *Node // 存储链表的最后一个节点内存地址
}

func (l LinkedList) Len() int {
	return l.Length
}
func (l LinkedList) Display(v int) *Node {
	// 遍历查找 n 的data值，并显示, 如果有 v 则返回
	n := l.Head /// 保持 Head头 元素在遍历的时候不会被改变
	for n != nil {
		fmt.Printf("%v -> ", n.Data)
		fmt.Println(l.Head.Data)
		if n.Data == v {
			return n
		}
		n = n.Next // 查看下一个
	}
	return nil
}

func (l *LinkedList) PushBack(n *Node) {
	// *LinkedList 如果类型没有 指针 *, 将导致L成为 没有实例化的对象
	// 推回 此方法将一个节点作为输入并将其链接到链表, 栈 式 压入
	if l.Head == nil {
		// 如果链表为空，则将传入的节点作为第一个节点，并且头尾都开始指向该节点，链表长度+1

		l.Head = n
		l.Tail = n
		l.Length += 1
	} else {
		// 当头节点存在时，执行else部分，尾节点下一个字段存储传入节点的内存地址，
		// 并且尾指向该节点

		l.Tail.Next = n
		l.Tail = n
		l.Length += 1
	}
}

/// 反转链表
func (nd *Node) reverse() *Node {
	// 空链表
	if nd == nil {
		return nil
	}
	var reverseHead *Node
	var preNode *Node
	curNode := nd
	for curNode != nil {
		nextNode := curNode.Next
		if nextNode == nil {
			reverseHead = curNode // 尾节点转为头节点
		}
		// 反转实现，当前节点前驱节点变为后驱节点
		curNode.Next = preNode
		// 设置下一个结点的前驱节点
		preNode = curNode
		curNode = nextNode
		/// 返回反转后的头结点
	}
	return reverseHead
}

// 遍历单链表
func (head *Node) traverse() {
	node := head
	for node != nil {
		fmt.Printf("%v ", node.Data)
		node = node.Next
	}
}

func main() {
	node := &Node{Data: 1} // 初始化，node变量将被传递给pushback方法，该方法将在链表中链接该节点
	fmt.Printf("%+v", node)

	ListLink := LinkedList{Length: 10, Head: node, Tail: &Node{Data: 10}}
	fmt.Printf("%+v", ListLink)
	fmt.Println(ListLink.Len())
	fmt.Printf("%+v\n", ListLink.Display(2))

	linked2 := &LinkedList{}
	for i := 0; i < 5; i++ {
		linked2.PushBack(&Node{Data: 11 + i})
		fmt.Printf("%+v", linked2)
		fmt.Println(linked2.Len())
		fmt.Printf("%+v\n", linked2.Display(2))
		// fmt.Println(Node)
	}

}
