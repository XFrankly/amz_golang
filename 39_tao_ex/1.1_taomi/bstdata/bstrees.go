package bstdata

import (
	"fmt"
)

/*
二叉搜索树是 红黑树的基础
红黑树 R-B tree， 全称 Red-Black Tree  本身是一个 二叉查找树，在其基础上附加了两个要求
树中每个热点增加一个用于存储颜色的 标志域
树中没有一个 路径比其他任何路径长出两倍，整个树要接近于 平衡 状态
*/

type TreeNode struct {
	Key           int
	Payload       string
	LeftChild     *TreeNode
	RightChild    *TreeNode
	Parent        *TreeNode
	balanceFactor int
}

type BinarySearchTree struct {
	Root *TreeNode
	Size int
}

// 缓存队列，用于存放 二叉树的 中序遍历结果
type CacheChan struct {
	Size  int /// cache 大小标记
	Read  <-chan *TreeNode
	Input chan<- *TreeNode
}

func NewCacheChan(size int) *CacheChan {
	if size <= 0 {
		panic("size must bigger than 0")
	}
	Chans := make(chan *TreeNode, size)
	return &CacheChan{
		Size:  size,
		Read:  Chans,
		Input: Chans,
	}
}

func (tn *TreeNode) HasLeftChild() *TreeNode {
	return tn.LeftChild
}

func (tn *TreeNode) HasRightChild() *TreeNode {
	return tn.RightChild
}

func (tn *TreeNode) IsLeftChild() bool {
	/// 是否左子 节点
	if tn.Parent != nil && tn.Parent.LeftChild == tn {
		return true
	}
	return false
}

func (tn *TreeNode) IsRightChild() bool {
	// 是否 右子节点
	if tn.Parent != nil && tn.Parent.RightChild == tn {
		return true
	}
	return false
}

func (tn *TreeNode) IsRoot() bool {
	// 是否根节点
	if tn.Parent == nil {
		return true
	} else {
		return false
	}
}

func (tn *TreeNode) IsLeaf() bool {
	// 是否叶子节点
	if tn.RightChild == nil && tn.LeftChild == nil {
		return true
	} else {
		return false
	}
}

func (tn *TreeNode) HasAnyChildren() *TreeNode {
	// 如果右子节点不为 空 返回右子节点，否则，
	// 查看左子节点，如果不为空，返回 左子节点，如果左右子节点都为nil，返回nil
	if tn.RightChild != nil {
		return tn.RightChild
	} else if tn.LeftChild != nil {
		return tn.LeftChild
	} else {
		return nil
	}
}

func (tn *TreeNode) HasBothChildren() bool {
	// 如果 两个子节点都存在，则返回 true，否则返回 false
	if tn.RightChild != nil && tn.LeftChild != nil {
		return true
	} else {
		return false
	}
}

// 遍历bst 树，查询 key 是否存在 该树中，如果存在，返回该节点，不存在，返回nil
func (tn *TreeNode) IterIsIn(key int) *TreeNode {
	if tn != nil {
		tnLeft := tn.HasLeftChild() //  左子树
		for tnLeft != nil {
			Logg.Printf("%+v\n", tnLeft)
			if tnLeft.Key == key {
				return tnLeft
			}
			if tnLeft == tn.LeftChild.LeftChild {
				break
			}
			tnLeft = tn.LeftChild.LeftChild
		}
		Logg.Printf("%+v\n", tn)
		if tn.Key == key {
			return tn
		}

		// defer MutilLock.Unlock()
		var tnRight *TreeNode
		tnRight = tn.HasRightChild() // 右子树
		for tnRight != nil {
			if tnRight.Key == key {
				return tnRight
			}

			// cMaps[tnRight.Key] = tnRight
			if tnRight == tn.RightChild.RightChild {
				break
			}
			tnRight = tn.RightChild.RightChild
			Logg.Printf("tnRight now:%+v cMaps:%v \n", tnRight, cMaps)
		}
	}
	return nil
}

// 创建一个 size 大小的 chan
func (tn *TreeNode) MakeCacheChan(size int) *CacheChan {
	return NewCacheChan(size)
}

// 向缓存 通道 存入 TreeNode 对象
func (tn *TreeNode) CachePuts(chans *CacheChan, newNode *TreeNode) *CacheChan {
	if len(chans.Input) < chans.Size {
		MutilLock.Lock()
		defer MutilLock.Unlock()
		chans.Input <- newNode
	} else {
		Logg.Println("its full cache channel.", chans)
	}
	return chans
}

func Doth() {
	bst2 := &BinarySearchTree{}
	for i := 0; i < 49; i++ {
		bst2.Puts(i, fmt.Sprintf("suanzi_%v", i))
		Display(bst2)
		// time.Sleep(time.Second * 1)
	}

	// Display(bst2)
	Logg.Printf("left:%v\n", bst2.Root.LeftChild)
	Logg.Printf("Right:%v\n", bst2.Root.RightChild)
	bst2.Rebalance(bst2.Root)
	//Display(bst2)

	caches := bst2.IterCache()
	Logg.Printf("caches:%v\n", caches)
	//查找某个节点 并以此旋转 再平衡
	keyNodes := bst2.Searcher(26)
	Logg.Printf("keyNodes:%v\n", keyNodes)

	//再平衡
	// bst2.RotateLeft(keyNodes)
	// Logg.Printf("bst2 :%v left:%v\n", bst2.Root, bst2.Root.LeftChild)

	bst2.Rebalance(keyNodes)
	Logg.Printf("bst2 left:%v right:%v\n", bst2.Root.LeftChild, bst2.Root.RightChild)

	//根据bst2 再造新树
	newTree := bst2.Rebuild(keyNodes)
	Logg.Printf("newTree :%v, root:%v\n", newTree, newTree.Root)

	//根节点作为 人才
	Logg.Printf("人 newTree.Root:%v\n", newTree.Root)
	//查看左子树 天才
	CacheChan := NewCacheChan(49)
	chans := IterCacheLeftNode(CacheChan, newTree.Root)
	Logg.Printf("天 chans :%#v\n", len(chans.Read))
	//查看右子树 地才
	CacheChanRight := NewCacheChan(49)
	chansRight := IterCacheRightNode(CacheChanRight, newTree.Root)
	Logg.Printf("地 chansRight:%#v\n", len(chansRight.Read))
}
