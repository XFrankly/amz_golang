package main

import (
	"log"
	"os"
)

/*
单环形链表
"""            -
           -       b
X---a---Y-            -
           c        Z
               -
给定一个链表，判断链表中是否有环。
    思路: 设定两个指针，fast 步长为2，slow 步长为1，遍历，
    从同一点出发，如果两个指针相遇，那么表示这个链表有环。
    如果没有环，那么fast 先达到 NULL

    快慢指针相遇时的循环次数等于环的长度。
进阶：
你能否不使用额外空间解决此题？

如何找出 环链表的入口
    思路: 快慢指针相遇点 到环入口的距离 = 链表起始点到环 入口的距离
如何求解环的长度
    思路:当两个指针相遇Z时，此时slow走过的路径为 a + b， fast走过的路径为 a+b+c+b, 因为fast步长为 slow的两边，所以 2（a+b）=a+b+c+b
    所以 a = c， 由图所示可知 环的长度为 b + c，也就是 a+b，即slow走过的步长
如何判断两个链表是否相交
    思路: 如果两个链表相交，那么这两个链表尾节点 一定相同，直接判断尾节点是否相同即可，
    环形链表可以拆分为两个Y链表。
"""
*/

var (
	Logg = log.New(os.Stderr, "INFO - ", 18)
)

type Node struct {
	Number int
	Prev   *Node
	Next   *Node
}

type Doublist struct {
	Lens int
	Head *Node
	Tail *Node
}

func MakeDlist() *Doublist {
	return &Doublist{}
}

func (th *Doublist) Lesseq(n *Node) (int, *Node) {
	/// 排序插入，找到相比 n 大于等于 的目标的结点
	if th.Lens <= 0 || n == nil {
		return 0, nil
	}
	currentNode := th.Head
	for i := 0; i < th.Lens; i++ {
		if currentNode.Number >= n.Number {
			return i, currentNode
		} else {
			currentNode = currentNode.Next
		}
	}
	/// 没有找到 比n 大的
	return th.Lens - 1, nil
}

/// 添加到空链表
func (th *Doublist) NewNodeList(n *Node) bool {
	////
	if th.Lens == 0 {
		th.Head = n
		th.Tail = n
		n.Prev = nil
		n.Next = nil
		th.Lens += 1
		return true
	} else {
		Logg.Panic("not empty node list.")
	}
	return false
}

// 头部添加 节点
func (th *Doublist) PushHead(n *Node) bool {
	//
	if th.Lens == 0 {
		return th.NewNodeList(n)
	} else {
		th.Head.Prev = n
		n.Prev = nil
		n.Next = th.Head
		th.Head = n
		th.Lens += 1
		return true
	}
}

//// 查找 第 index 节点
func (th *Doublist) Get(ind int) *Node {
	currentNode := th.Head
	for i := 0; i < th.Lens; i++ {
		if i == ind {
			return currentNode
		} else {
			currentNode = currentNode.Next
		}
	}
	return nil
}

//// 摘取 第index节点
func (th *Doublist) Pick(ind int) *Node {
	currentNode := th.Get(ind)
	th.Remove(currentNode)
	return currentNode
}

func (th *Doublist) Remove(n *Node) {
	/// index
	if n == nil || n == th.Tail {
		n = th.Tail
		th.Tail = n.Prev
		th.Tail.Next = nil
	} else if n == th.Head {
		th.Head = n.Next
		th.Head.Prev = nil
	} else {
		n.Prev.Next = n.Next
		n.Next.Prev = n.Prev
	}
	th.Lens -= 1
}

/// 从头部 获取一个 节点
func (th *Doublist) PopHead() *Node {
	if th.Lens <= 0 {
		return nil
	}
	cNode := th.Pick(0)
	return cNode
}

//// 从尾部获取一个 节点
func (th *Doublist) PopTail() *Node {
	if th.Lens <= 0 {
		return nil
	}
	cNode := th.Pick(th.Lens - 1)
	return cNode
}

/// 添加尾部节点
func (th *Doublist) Append(n *Node) bool {
	///
	if th.Lens == 0 {
		return th.NewNodeList(n)
	} else {
		th.Tail.Next = n
		n.Prev = th.Tail
		n.Next = nil
		th.Tail = n
		th.Lens += 1
		return true
	}
}

func (th *Doublist) Pushback(n *Node) bool {
	///
	if n == nil {
		return false
	}
	currentNode := th.Head
	if currentNode == nil {

		return th.NewNodeList(n)
	} else {
		inDex, insertNode := th.Lesseq(n)
		if inDex == 0 {
			return th.PushHead(n)
		} else if inDex == (th.Lens-1) && insertNode == nil {
			return th.Append(n)
		}
		Logg.Printf("insert at :%+v\n", inDex)

		n.Next = insertNode
		n.Prev = insertNode.Prev
		//// 很容易失败
		if insertNode.Prev != nil {
			insertNode.Prev.Next = n
		}

		insertNode.Prev = n
		th.Lens += 1
		return true
	}
}

func (th *Doublist) Display() {
	/// 显示链表的值
	node := th.Head
	t := 0
	Logg.Println(node.Number)
	for node != nil {
		t += 1
		if t >= th.Lens {
			break
		}

		node = node.Next
		Logg.Println(node.Number)
	}

	Logg.Println("length:", th.Lens)
}

////////////////////////////////// 基于双向链表的 环形链表
/*
环形链表用例:  pos 表示 链表尾部节点 执行该链表的 哪一个节点
		输入:
		[3,2,0,-4]    pos = 1  # 表示1个环 -4 与 2相 链
			|____|

		输出:
		tail connects to node index 1

		输入:
			[1, 3, 4]  pos =0   #  表示 4 与 1 首尾相连
			|_____|
		输出:
			tail connects to node index 0

		输入:
			[1, 2]  pos = 0  # 表示 2 与 1 连
			|__|
		输出:
			tail connects to node index 0

		输入:
			[-1]
			-1  # 表示没有 环
		输出:
			null
*/
type CycleLinked struct {
	Size  int
	Items *Doublist
}
type MyChan struct {
	Read  <-chan *Node
	Input chan<- *Node
}

/// Dlist 双向链表，pos 表示 该双向链表的尾节点 与 该链表哪一个节点相连
func MakeCycleLinked(Dlist *Doublist, pos int) *CycleLinked {
	NodePos := Dlist.Get(pos)
	Dlist.Tail.Next = NodePos
	Dlist.Lens += 1 /// 保证取到 环形链表的最后一个节点 pos
	return &CycleLinked{Size: Dlist.Lens, Items: Dlist}
}

func MakeEmptyLinked(n int) *CycleLinked {
	return &CycleLinked{Size: n}
}
func (cl *CycleLinked) GetItems() *Doublist {
	return cl.Items
}

/// 添加n个 节点到 链表
func (cl *CycleLinked) SetItems(n int) *CycleLinked {
	if n <= cl.Size {
		cl.Items = MakeDlist()
		for i := 0; i < n; i++ {
			cl.Items.Pushback(&Node{Number: i})
		}
	}

	return cl

}

/// 返回一个从 第0个 开始的单环链表
func (cl *CycleLinked) MakeChanCycle() *MyChan {
	/// 将链表取出 放入 通道
	var noe *Node
	var MyC = make(chan *Node, cl.Size)
	var chans = &MyChan{Read: MyC, Input: MyC}
	for i := 0; i < cl.Items.Lens; i++ {
		noe = cl.Items.PopHead()
		cl.Items.Append(noe)
		chans.Input <- noe
	}
	return chans
}

func DoubleLinkedList() *Doublist {
	dlist := MakeDlist()
	slit := []int{9, 2, 5, 6, 7, 2, 6, 10, 3}
	for _, i := range slit {
		node := &Node{Number: i}
		dlist.Pushback(node)
	}

	dlist.Display()
	Logg.Println()
	dlist.Pushback(&Node{Number: 123})
	dlist.Display()

	Logg.Println()
	dlist.PopHead()
	dlist.Display()
	Logg.Println()

	dlist.PopTail()
	dlist.Display()
	Logg.Println()

	dlist.Pick(2)
	dlist.Display()
	Logg.Println()

	dlist.Pushback(&Node{Number: 0})
	dlist.Display()
	Logg.Println()
	dlist.Pushback(&Node{Number: 1})
	dlist.Display()
	Logg.Println()

	return dlist
}

func main() {
	dlist := DoubleLinkedList()

	clist := MakeCycleLinked(dlist, 2) /// 尾节点与 第2个节点链接
	cChans := clist.MakeChanCycle()
	Logg.Printf("%+v\n", cChans)
	for len(cChans.Read) > 0 {
		cr := <-cChans.Read
		Logg.Printf("%d, %+v\n", cr.Number, cr)
	}

	///// 空链表
	tLinked := MakeEmptyLinked(8)
	//// 填入数字
	numberC := tLinked.SetItems(tLinked.Size)

	cNode := numberC.Items.Head
	for cNode != nil {
		Logg.Printf("%+v\n", cNode)
		cNode = cNode.Next
	}
}
