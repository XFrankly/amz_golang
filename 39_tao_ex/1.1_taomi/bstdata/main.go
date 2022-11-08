package bstdata

import (
	"log"
	"os"
	"sync"
)

var (
	cMaps     = make(map[int]*TreeNode)
	Logg      = log.New(os.Stderr, "INFO -:", 13)
	MutilLock sync.Mutex
	WG        sync.WaitGroup
)

// 遍历节点的 右子树
func IterCacheRightNode(ccChan *CacheChan, tnode *TreeNode) *CacheChan {
	tnRight := tnode.HasRightChild() // 右子树
	for tnRight != nil {
		tnode.CachePuts(ccChan, tnRight)
		Logg.Printf("%+v\n", tnRight)
		/// 右子节点的 左子节点遍历
		if tnRight.HasLeftChild() != nil {
			ccChan = IterCacheLeftNode(ccChan, tnRight)
		}
		tnRight = tnRight.RightChild

	}
	return ccChan
}

// 遍历节点的 左子树
func IterCacheLeftNode(ccChan *CacheChan, tnode *TreeNode) *CacheChan {
	tnLeft := tnode.HasLeftChild() // 右子树
	for tnLeft != nil {
		tnode.CachePuts(ccChan, tnLeft)
		Logg.Printf("%+v\n", tnLeft)
		/// 左子节点的 右子节点遍历
		if tnLeft.HasRightChild() != nil {
			ccChan = IterCacheRightNode(ccChan, tnLeft)
		}
		tnLeft = tnLeft.LeftChild

	}
	return ccChan
}
