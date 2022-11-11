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

// 添加一个节点到缓存
func (ch *CacheChan) Putin(td *TreeNode) int {
	MutilLock.Lock()
	defer MutilLock.Unlock()

	if len(ch.Input) < ch.Size {
		ch.Input <- td
		return ch.Size
	}
	return 0
}

// 按序取出一个node 节点
func (ch *CacheChan) GetNode() *TreeNode {
	MutilLock.Lock()
	defer MutilLock.Unlock()
	if len(ch.Read) > 0 {
		td := <-ch.Read
		return td
	}
	return nil

}

func (ch *CacheChan) IsTreeNodeIn(key int) bool {
	isIn := false
	if ch.Size <= 0 {
		return isIn
	}
	t := ch.Size
	for i := 0; i < t; i++ {
		td := ch.GetNode()

		ch.Putin(td)
		if td.Key == key {
			isIn = true
		}
	}
	return isIn
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
			if tnLeft == tn.LeftChild {
				break
			}
			tnLeft = tn.LeftChild
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
			if tnRight == tn.RightChild {
				break
			}
			tnRight = tn.RightChild
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
		RLock.Lock()
		defer RLock.Unlock()
		chans.Putin(newNode)
	} else {
		Logg.Println("its full cache channel.", chans)
	}
	return chans
}

// 中序遍历
// func (tn *TreeNode) NodeKeyList(size int) []int {
// 	var tnd []int
// 	if tn != nil { // # 根不为空，基本结束条件
// 		tnLeft := tn.HasLeftChild()
// 		if tnLeft != nil { //:  # 左子树不为空
// 			for tn.LeftChild != nil {
// 				tnd = append(tnd, tn.LeftChild.Key)
// 				//# 遍历左子树
// 				//yield elem  //# 生成器，返回左子树一个元素
// 				tn.LeftChild = tn.LeftChild.LeftChild
// 			}
// 		}
// 		tnd = append(tnd, tn.Key) //yield self.key  # 生成器，返回根
// 		tnRight := tn.HasRightChild()
// 		if tnRight != nil { //:  # 右子树不为空
// 			//for elem in self.rightChild //:  # 遍历右子树
// 			//yield elem  # 生成器，返回右子树一个元素
// 			for tn.RightChild != nil {
// 				tnd = append(tnd, tn.RightChild.Key)
// 				tn.RightChild = tn.RightChild.RightChild
// 			}
// 		}
// 	}
// 	return tnd
// }

// 中序遍历 二叉树节点的子节点 并存储到channel 返回指向 channel的指针
func (tn *TreeNode) IterCache(size int) *CacheChan {
	cacheCahn := tn.MakeCacheChan(size)
	if tn != nil {
		tnLeft := tn.HasLeftChild() //  左子树
		if tnLeft != nil {
			cacheCahn.Putin(tnLeft)
			cacheCahn = IterCacheLeftNode(cacheCahn, tnLeft)
		}
		// for tnLeft != nil { //// 遍历左子树的左节点
		// 	Logg.Printf("%+v\n", tnLeft)
		// 	cacheCahn.Putin(tnLeft)
		// 	if tnLeft == tnLeft.LeftChild {
		// 		Logg.Printf("left child same:%v\n", tnLeft)
		// 		break
		// 	}
		// 	tnLeft = tnLeft.LeftChild
		// }
		Logg.Printf("%+v\n", tn)
		cacheCahn.Putin(tn)
		tnRight := tn.HasRightChild() // 右子树
		if tnRight != nil {
			cacheCahn.Putin(tnRight)
			cacheCahn = IterCacheRightNode(cacheCahn, tnRight)
		}
		// for tnRight != nil {
		// 	tn.CachePuts(cacheCahn, tnRight)
		// 	Logg.Printf("%+v\n", tnRight)
		// 	if tnRight == tnRight.RightChild {
		// 		Logg.Printf("right child same:%v\n", tnRight)
		// 		break
		// 	}
		// 	tnRight = tnRight.RightChild
		// }
	}
	return cacheCahn
}

// 当前节点的左子节点，左子树的 最左下角的 值
func (tn *TreeNode) FindMin() *TreeNode {
	current := tn                       // 根节点
	for current.HasLeftChild() != nil { // 直到找到最左下角的值，就是直接后继
		current = current.LeftChild
	}
	return current
}

// 当前节点的右子节点，右子树的 最右下角的 值
func (tn *TreeNode) FindMax() *TreeNode {
	current := tn                        // 根节点
	for current.HasRightChild() != nil { // 直到找到最右下角的值，就是直接后继
		current = current.RightChild
	}
	return current
}

// 调整平衡树
func (tn *TreeNode) ReplaceNodeData(key int, value string, lc *TreeNode, rc *TreeNode) {
	tn.Key = key
	tn.Payload = value
	tn.LeftChild = lc
	tn.RightChild = rc
	if tn.HasLeftChild() != nil {
		tn.LeftChild.Parent = tn
	}

	if tn.HasRightChild() != nil {
		tn.RightChild.Parent = tn
	}
}

// 摘出某个节点
func (tn *TreeNode) SpliceOut() {
	if tn.IsLeaf() {
		// 摘出叶子节点
		if tn.IsLeftChild() {
			tn.Parent.LeftChild = nil
		} else {
			tn.Parent.RightChild = nil
		}
	} else if tn.HasAnyChildren() != nil {
		if tn.HasLeftChild() != nil { // 摘 左子节点
			if tn.IsLeftChild() {
				// 这一代码块 在同时有两个左右子树，有左下子树的情况，不会执行该代码
				tn.Parent.LeftChild = tn.LeftChild
			} else {
				tn.Parent.RightChild = tn.LeftChild
			}
		} else {
			// 摘 右子节点
			if tn.IsLeftChild() {
				tn.Parent.LeftChild = tn.RightChild
			} else {
				// 摘 带右子节点的节点
				tn.Parent.RightChild = tn.RightChild
			}
			tn.RightChild.Parent = tn.Parent
		}
	}
}

// 寻找后继节点
func (tn *TreeNode) FindSuccessor() *TreeNode {
	var succe *TreeNode
	if tn.HasRightChild() != nil {
		succe = tn.RightChild.FindMin() /// 左子节点 直接后继
	} else {
		if tn.Parent != nil {
			/// 该节点没有 右子树，需要去其他地方找后继
			/// 在本例中，前提就是当前节点同时有 左右子树
			if tn.IsLeftChild() {
				succe = tn.Parent
			} else {
				tn.Parent.RightChild = nil
				succe = tn.Parent.FindSuccessor() // 递归调用查找
				tn.Parent.RightChild = tn
			}
		}
	}
	return succe
}

// 树管理
type BinarySearchTree struct {
	Root      *TreeNode //根节点
	Size      int       //本树总大小 包括根节点
	LeftSize  int       //左子树大小
	RightSize int       //右子树大小
}

// 搜索树 中序 遍历后的全部节点 队列
func (bst *BinarySearchTree) IterCache() *CacheChan {

	cache := bst.Root.IterCache(bst.Size)
	return cache
}

// 从缓存队列获取一个节点，如果有
func (bst *BinarySearchTree) CacheGets(chans *CacheChan) *TreeNode {

	if len(chans.Read) > 0 {
		newReader := <-chans.Read
		return newReader
	}
	return nil
}

// 插入节点
func (bst *BinarySearchTree) Put(key int, val string, cNode *TreeNode) {
	if key < cNode.Key {
		// 如果参数key比当前节点key 小，进入树的左子树进行递归插入
		if cNode.HasLeftChild() != nil {
			bst.Put(key, val, cNode.LeftChild) /// 递归左子树 插入
		} else {
			cNode.LeftChild = &TreeNode{Key: key,
				Payload: val, Parent: cNode} //树的左子节点
		}
	} else { /// 如果参数 key的值 大于当前节点key，进入树的右子树进入递归插入
		if cNode.HasRightChild() != nil {
			bst.Put(key, val, cNode.RightChild) /// 递归右子树
		} else {
			cNode.RightChild = &TreeNode{Key: key,
				Payload: val, Parent: cNode}
		}
	}
}

//	高度log2_n,如果key 列表随机分布，大于小于根节点的key的键值 大致相当
//
// 性能在于二叉树的高度，最大层次，高度也受数据项key插入顺序影响
// 算法复杂度 最差 O(log2_n)
func (bst *BinarySearchTree) Puts(key int, val string) bool {
	if bst.Root != nil {
		// 有根节点
		MutilLock.Lock()
		defer MutilLock.Unlock()

		if bst.Searcher(key) != nil {
			// 已经存在 无法插入
			msg := fmt.Sprintf("the key had exist at treenode:%+v\n", key)
			Logg.Println(msg)
			return false
		}
		bst.Put(key, val, bst.Root)
		//更新大小信息
		if bst.Root.Key < key {
			//新 key 添加在右侧
			bst.RightSize += 1
		} else if bst.Root.Key >= key {
			//新 key 添加在右侧
			bst.LeftSize += 1
		}

	} else {
		/// 没有根节点
		bst.Root = &TreeNode{Key: key, Payload: val}
	}
	bst.Size += 1
	return true
}

// 设置节点
func (bst *BinarySearchTree) SetNode(node *TreeNode) bool {
	return bst.Puts(node.Key, node.Payload)

}

// 找到节点为key的 Payload值，只要是平衡树，get的时间复杂度可用保持在 O(logN)
func (bst *BinarySearchTree) Searcher(key int) *TreeNode {
	if bst.Root != nil {
		res := bst.Get(key, bst.Root) /// 递归该树
		if res != nil {
			return res
		}
	}
	return nil
}

// 当前节点，即要插入的 二叉查找树， 子树的根，为当前节点
func (bst *BinarySearchTree) Get(key int, cNode *TreeNode) *TreeNode {
	if cNode == nil {
		return nil
	} else if cNode.Key == key {
		return cNode
	} else if key < cNode.Key {
		return bst.Get(key, cNode.LeftChild)
	} else {
		return bst.Get(key, cNode.RightChild)
	}
}

// delete 的具体实现，要求仍然保持BST 性质
// 1 节点无子节点 2 节点有1个子节点 3 节点有2个子节点
func (bst *BinarySearchTree) remove(cNode *TreeNode) {
	if cNode.IsLeaf() {
		/// leaf 叶子节点，没有子节点，属于场景1，无子节点，直接删除
		if cNode == cNode.Parent.LeftChild {
			/// 本身是 左子节点
			cNode.Parent.LeftChild = nil
		} else {
			cNode.Parent.RightChild = nil
		}
	} else if cNode.HasBothChildren() {
		/// 有两个子节点
		succe := cNode.FindSuccessor() // 找到当前需要删除的节点的后继节点
		succe.SpliceOut()
		cNode.Key = succe.Key         // 替换Key
		cNode.Payload = succe.Payload // 替换Payload 值，节点的数据
	} else {
		/// 有一个子节点
		if cNode.HasLeftChild() != nil {
			if cNode.IsLeftChild() {
				/// 左子节点删除
				cNode.LeftChild.Parent = cNode.Parent    // 修改指针。当前节点的左子节点的父节点，修改为节点的父节点
				cNode.Parent.LeftChild = cNode.LeftChild // 修改指针，当前节点的父节点的左子节点，修改为当前节点的左子节点
			} else if cNode.IsRightChild() {
				/// 右 子节点删除
				cNode.LeftChild.Parent = cNode.Parent
				cNode.Parent.RightChild = cNode.LeftChild
			} else {
				// 根节点删除
				cNode.ReplaceNodeData(
					cNode.LeftChild.Key,
					cNode.LeftChild.Payload,
					cNode.LeftChild.LeftChild,
					cNode.LeftChild.RightChild,
				)
			}
		} else {
			if cNode.IsLeftChild() {
				/// 左子节点删除
				cNode.RightChild.Parent = cNode.Parent
				cNode.Parent.LeftChild = cNode.RightChild
			} else if cNode.IsRightChild() {
				/// 右子节点删除
				cNode.RightChild.Parent = cNode.Parent
				cNode.Parent.RightChild = cNode.RightChild
			} else {
				/// 根节点删除
				cNode.ReplaceNodeData(
					cNode.RightChild.Key,
					cNode.RightChild.Payload,
					cNode.RightChild.LeftChild,
					cNode.RightChild.RightChild,
				)
			}
		}
	}
}

// deletes 用于删除 树中某个节点，子节点替换当前节点，具体是调用 delete方法
// 修改树的 大小信息
func (bst *BinarySearchTree) Deletes(key int) {
	var nTRemove *TreeNode
	if bst.Size > 1 {
		nTRemove = bst.Get(key, bst.Root)
		if nTRemove != nil {
			bst.remove(nTRemove)
			bst.Size -= 1
		} else {
			msg := "Error, key not in tree"
			panic(msg)
		}
	} else if bst.Size == 1 && bst.Root.Key == key {
		nTRemove = bst.Root
		bst.Root = nil
		bst.Size -= 1
	} else {
		msg := "Error, key not in tree."
		panic(msg)
	}
	//没有报错 则更新大小
	if bst.Root.Key > key {
		//此key在左子树
		bst.LeftSize -= 1
		Logg.Printf("success del nTRemove:%v bst.LeftSize :%v \n", nTRemove, bst.LeftSize)
	} else if bst.Root.Key <= key {
		// 此key 在右子树
		bst.RightSize -= 1
		Logg.Printf("success del nTRemove:%v bst.RightSize :%v \n", nTRemove, bst.RightSize)

	}

}

/*
@param side  在哪一侧执行删除，0 左侧,1 右侧
@param numb 删除多少 节点，如果需要删除的节点树 大于 该侧树大小，返回false
*/
//// 测试bst平衡树
func (bst *BinarySearchTree) YaoDel(side, numb int) (bool, map[string][]*TreeNode) {
	var nodeDels = make(map[string][]*TreeNode, 27)
	if bst.Size < numb {
		return false, nil
	}

	if side == 0 { //左侧
		if bst.LeftSize < numb {
			return false, nil
		}
		var leftDelNodes = []*TreeNode{}
		for numb > 0 {
			//找到要删除的最小节点,根节点出发
			mt := bst.Root.FindMin()
			//删除 numb 个左子树节点
			bst.Deletes(mt.Key)
			Logg.Printf("after del :%v left 1 key left size:%v, right size:%v\n", mt, bst.LeftSize, bst.RightSize)
			numb -= 1
			leftDelNodes = append(leftDelNodes, mt)
		}
		nodeDels["left"] = leftDelNodes
	}

	if side == 1 { //右侧执行删除
		if bst.RightSize < numb {
			return false, nil
		}
		var rightDelNodes = []*TreeNode{}
		for numb > 0 {
			//找到要删除的值最大的节点,根节点出发
			mx := bst.Root.FindMax()
			//删除 numb 个右子树节点
			bst.Deletes(mx.Key)
			Logg.Printf("after del :%v RightChild 1 key left size:%v, right size:%v\n", mx, bst.LeftSize, bst.RightSize)
			numb -= 1
			rightDelNodes = append(rightDelNodes, mx)
		}
		nodeDels["right"] = rightDelNodes
	}

	return true, nodeDels
}

// 更新平衡树
func (bst *BinarySearchTree) UpdateBalance(tn *TreeNode) {
	if tn.balanceFactor > 1 || tn.balanceFactor < -1 {
		bst.Rebalance(tn) // 重新平衡
	}

	if tn.Parent != nil {
		/// 查看当前节点是否 有父节点，如果没有，说明是根节点，无需再传播
		if tn.IsLeftChild() {
			tn.Parent.balanceFactor += 1
		} else if tn.IsRightChild() {
			tn.Parent.balanceFactor -= 1
		}
		if tn.Parent.balanceFactor != 0 {
			/// 如果父节点平衡因子不为0，进行父节点平衡因子的调整
			bst.UpdateBalance(tn.Parent) // 调整父节点因子
		}
	}
}

// 以节点 tn 再造平衡树, 根据制定节点 重新生成一个平衡子树, 必须是有序连续的
func (bst *BinarySearchTree) Rebuild(tn *TreeNode) *BinarySearchTree {
	bstNew := &BinarySearchTree{}
	bstNew.SetNode(tn) //根节点
	// cChans := bst.IterCache() //中序遍历 无法完整遍历，？？？
	min := bst.Root.FindMin().Key
	max := bst.Root.FindMax().Key
	Logg.Printf("min:%v max:%v max newSize:%v \n", min, max, max-min+1)
	for min < max+1 { //, m :=  cMaps.Size {
		Tnode := &TreeNode{Key: min, Payload: fmt.Sprintf("suanzi_%v", min)} //cChans.GetNode() //bst.CacheGets(cChans)
		if Tnode == nil || Tnode.Key == tn.Key {
			min += 1
			continue
		}
		bstNew.SetNode(Tnode)
		min += 1
	}
	Logg.Printf("new tree size:%v \n", bstNew.Size)
	return bstNew
}

// 节点子树再平衡，左或右旋转
func (bst *BinarySearchTree) Rebalance(tn *TreeNode) {
	if tn.balanceFactor < 0 {
		// 右子树 重，需要旋转
		if tn.RightChild.balanceFactor > 0 {
			//  做一个 LR 旋转， LR Rotation
			bst.RotateRight(tn.RightChild) /// 右子节点 左重，先右旋
			bst.RotateLeft(tn)
		} else {
			/// 单次 左旋
			bst.RotateLeft(tn)
		}
	} else if tn.balanceFactor > 0 {
		if tn.LeftChild.balanceFactor < 0 {
			/// 左重需要右旋
			bst.RotateLeft(tn.LeftChild) /// 左子节点右重 先左 旋转
			bst.RotateRight(tn)
		} else {
			/// 单次右旋
			bst.RotateRight(tn)
		}
	}
}

// 在指定节点tn处，旋转左子树，旋转调整左子树平衡
func (bst *BinarySearchTree) RotateLeft(tn *TreeNode) {
	newRoot := tn.RightChild
	if newRoot.LeftChild != nil {
		tn.RightChild = newRoot.LeftChild
		newRoot.LeftChild.Parent = tn
	}
	newRoot.Parent = tn.Parent
	if tn.IsRoot() {
		bst.Root = newRoot
	} else {
		if tn.IsLeftChild() {
			tn.Parent.LeftChild = newRoot
		} else {
			tn.Parent.RightChild = newRoot
		}
	}
	newRoot.LeftChild = tn
	tn.Parent = newRoot

	/// 仅有两个节点需要调整因子
	lessBf := 0
	moreBf := 0
	if newRoot.balanceFactor < 0 {
		lessBf = newRoot.balanceFactor
	} else {
		moreBf = newRoot.balanceFactor
	}
	tn.balanceFactor = tn.balanceFactor + 1 - lessBf
	newRoot.balanceFactor = newRoot.balanceFactor + 1 + moreBf
}

// // 在指定节点tn处，右旋转，调整右子树平衡
func (bst *BinarySearchTree) RotateRight(tn *TreeNode) {
	newRoot := tn.LeftChild
	if newRoot.RightChild != nil {

		tn.LeftChild = newRoot.RightChild
		newRoot.RightChild.Parent = tn
	}

	newRoot.Parent = tn.Parent
	if tn.IsRoot() { //
		bst.Root = newRoot
	} else {
		if tn.IsRightChild() {
			tn.Parent.RightChild = newRoot
		} else {
			tn.Parent.LeftChild = newRoot
		}
	}
	newRoot.RightChild = tn
	tn.Parent = newRoot
	/// 仅有两个节点需要调整因子
	lessBf := 0
	moreBf := 0
	if newRoot.balanceFactor < 0 {
		lessBf = newRoot.balanceFactor
	} else {
		moreBf = newRoot.balanceFactor
	}
	tn.balanceFactor = tn.balanceFactor + 1 - lessBf
	newRoot.balanceFactor = newRoot.balanceFactor + 1 + moreBf
}
