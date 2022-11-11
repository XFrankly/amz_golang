package bstdata

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	cMaps       = make(map[int]*TreeNode)
	BstNew      = MakeNewTreeManager()
	Logg        = log.New(os.Stderr, "INFO -:", 13)
	MutilLock   = sync.RWMutex{}
	RLock       = sync.RWMutex{}
	WG          sync.WaitGroup
	defMean     = map[int]string{6: "老阴", 7: "少阳", 8: "少阴", 9: "老阳"}
	Coordinates = map[int]string{1: "初", 2: "二", 3: "三", 4: "四", 5: "五", 6: "上"}
	EnvRange    = map[string]int{"A": -1, "B": 0, "C": 1}
	start       = `
           _______
           _______
           _______
           .......
___ ___....元亨利贞...._______
___ ___....运势占卜...._______
_______....大吉大利....___ ___ 
           .......
           ___ ___
           ___ ___
           ___ ___

`
)

// 遍历节点的 右子树
func IterCacheRightNode(ccChan *CacheChan, tnode *TreeNode) *CacheChan {

	tnRight := tnode.HasRightChild() // 右子树
	for tnRight != nil {
		// tnode.CachePuts(ccChan, tnRight)
		ccChan.Putin(tnRight)
		Logg.Printf("add right child node:%+v\n", tnRight)
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
		ccChan.Putin(tnLeft)
		Logg.Printf("add left child node: %+v\n", tnLeft)
		/// 左子节点的 右子节点遍历
		if tnLeft.HasRightChild() != nil {
			ccChan = IterCacheRightNode(ccChan, tnLeft)
		}
		tnLeft = tnLeft.LeftChild

	}
	return ccChan
}

// 读取并消费所有节点并 有序存入双端链表
func IterNodes(ccChan *CacheChan) *dlist {

	ccList := &dlist{}
	if ccChan.Size <= 0 {
		msg := fmt.Sprintf("ccChan size error:%v\n", ccChan)
		panic(msg)
	}

	cr := len(ccChan.Read)
	Logg.Printf("chan total size:%v\n", cr)
	for len(ccChan.Read) > 0 {
		nowNode := ccChan.GetNode()
		if nowNode == nil {
			continue
		}
		newTNode := &node{numb: nowNode}
		if _, node := ccList.isNodeIn(newTNode); node != nil {
			continue
		}

		ccList.pushback(newTNode)
	}

	return ccList
}

type TaoBstManager struct {
	*BinarySearchTree
	Bains       []int
	Indent      int
	DefMean     map[int]string
	Coordinates map[int]string
}

func MakeNewTreeManager() *TaoBstManager {

	return &TaoBstManager{
		BinarySearchTree: &BinarySearchTree{},
		Indent:           8,
		DefMean:          defMean,
		Coordinates:      Coordinates,
	}
}

func (bsts *TaoBstManager) Trees() *BinarySearchTree {
	return bsts.BinarySearchTree
}

// 遍历树所有节点 到channel
func (bsts *TaoBstManager) IterListNodes() *dlist {

	ccHan := bsts.IterCache()
	Logg.Printf("ccHan size:%v\n", ccHan.Size)
	return IterNodes(ccHan)
}

// 显示节点
func (bsts *TaoBstManager) Display() {

	n := 0
	Logg.Printf("%+v\n", bsts)
	Logg.Println("only root tree:", bsts.Root, bsts.Size)
	cChans := bsts.IterCache()
	Logg.Printf("Display ccHan size:%v\n", cChans.Size)
	for i := 0; i < cChans.Size; i++ { //, m :=  cMaps.Size {
		Tnode := cChans.GetNode() // bsts.CacheGets(cChans)
		Logg.Println("map index:", Tnode)
		n += 1
		if Tnode == nil {
			Logg.Println("had show all bst node size:", cChans.Size, "after index is nil:", i)
			continue
		}
		Logg.Printf("the index:%+v\n", i)
	}
	Logg.Println("had after show all bst node size:", len(cChans.Read), "total show:", n)
}

// 做一变
func (newTree *TaoBstManager) DoOneBian() (*TaoBstManager, int) {

	//一变
	lys := newTree.LeftSize % 4
	if lys == 0 {
		lys = 4
	}
	ldelRst, leftNodeDel := newTree.YaoDel(0, lys) //当左侧长度 除以4的 余数为3时，从左侧删除3个节点

	rys := newTree.RightSize % 4
	if rys == 0 {
		rys = 4
	}
	rdelRst, rightNodeDel := newTree.YaoDel(1, rys+1) //从右节点多删除一个，当作人材的归档，当左侧长度 除以4的 余数为3时，从左侧删除3个节点

	Logg.Printf("del left lys:%v, leftNodeDel:%v del right rys:%v rightNodeDel:%v,\n", lys, leftNodeDel, rys, rightNodeDel)
	Logg.Printf("DoOneBian letf:%v size:%v, right:%v size:%v \n", ldelRst, newTree.LeftSize, rdelRst, newTree.RightSize)
	return newTree, len(leftNodeDel) + len(rightNodeDel)
}

// 为一爻做一变
func (newManager *TaoBstManager) DoOneBianForYao(env string) *TaoBstManager {

	// 计算随机数，必须在大于4 ~ Max-4 范围内，否则一变失败
	tm, _ := newManager.DoOneBian()
	newTree := tm.Trees()
	mi := newTree.Root.FindMin() //最小
	mx := newTree.Root.FindMax() //最大
	//st := newTree.Size           //总数
	//新的根节点 和以此生成的新的树 作为二变基础
	two := DoRange(mi.Key, mx.Key, env)

	keyNodesTwo := newTree.Searcher(two)
	Logg.Printf("new tree mi:%v mx:%v root:%v node:%v\n", mi, mx, two, keyNodesTwo)
	newTreeTwo := newTree.Rebuild(keyNodesTwo)

	nms := &TaoBstManager{BinarySearchTree: newTreeTwo}

	return nms

}

// 做一爻 三变
func (newManager *TaoBstManager) DoOneYao(env string) int {

	// 计算随机数，必须在大于4 ~ Max-4 范围内，否则一变失败
	newMan := newManager.DoOneBianForYao(env)
	Logg.Printf("newMan of:%v\n", newMan.Size)
	newTwoMan := newMan.DoOneBianForYao(env)
	Logg.Printf("newTwoMan of:%v\n", newTwoMan.Size)

	newThreeMan := newTwoMan.DoOneBianForYao(env)
	Logg.Printf("newThreeMan of:%v, yao:%v, \n whichs deleted:%v \n", newThreeMan.Size, newThreeMan.Size/4, newThreeMan.Bains)

	return newThreeMan.Size / 4
}

// 初始化树管理器 并创建六爻
func (bstNew *TaoBstManager) YaoWithReBuild(t int, env string) ([]int, *TaoBstManager) {

	if t <= 0 {
		t = 49
	}
	for i := 0; i < t; i++ {
		bstNew.Puts(i, fmt.Sprintf("suanzi_%v", i))
		Logg.Printf("bstNew :%v, root:%v  len:%v\n  \n", bstNew, bstNew.Root, bstNew.Size)

	}

	var yaoOnes []int

	//得一卦象
	for len(yaoOnes) < 6 {
		// 以人为本 分天地人 查找某个节点 并以此旋转 再平衡
		rootKey := DoRand(t)
		keyNodes := bstNew.Searcher(rootKey)
		//根据bst 再造新树
		newTree := bstNew.Rebuild(keyNodes)
		caches := newTree.IterCache()
		Logg.Printf("newTree :%v, root:%v caches len:%v\n cachesInt:%v\n", newTree.Size, newTree.Root, caches.Size, len(caches.Read))
		Logg.Printf("树大小没有变:%v, new:%v, equal:%v\n", bstNew.Size, newTree.Size, bstNew.Size == newTree.Size)
		if bstNew.Size != newTree.Size {
			msg := fmt.Sprintf("size error with new tree:%v\n", newTree.Size)
			panic(msg)
		}
		treeManager := &TaoBstManager{BinarySearchTree: newTree}
		yaoOne := treeManager.DoOneYao(env)
		yaoOnes = append(yaoOnes, yaoOne)
	}

	Logg.Printf("yaoOnes:%v\n", yaoOnes)
	return yaoOnes, bstNew
}

// 显示和解释,输出二进制格式 和变爻位置
func (dt *TaoBstManager) AnysisYaos(yaos []int) ([]int, map[int]bool) {
	//是否为 变卦后的
	cp := make(map[int]bool, 6)
	//转换为 二进位
	newYaos := []int{}

	for i, g := range yaos {
		if g == 6 || g == 8 {
			if g == 6 {
				dt.FormatShow("__ __ (6 " + dt.DefMean[6] + ")\n")
				cp[i] = true //# 变卦 为 true 只写入变卦后的 少阳 少阴
				//老阴变少阳
				newYaos = append(newYaos, 1)
			}
			if g == 8 {
				dt.FormatShow("__ __ (8 " + dt.DefMean[8] + ")\n")
				cp[i] = false                //# 只写入变卦后的 少阳 少阴
				newYaos = append(newYaos, 0) //# 只写入变卦后的 少阳 少阴
			}

		} else if g == 7 || g == 9 {
			if g == 7 {
				dt.FormatShow("_____ (7 " + dt.DefMean[7] + ")\n")
				cp[i] = false
				newYaos = append(newYaos, 1)
			}

			if g == 9 {
				dt.FormatShow("_____ (9 " + dt.DefMean[9] + ")\n")
				//老阳变少阴
				cp[i] = true
				newYaos = append(newYaos, 0)
			}

		} else {
			dt.FormatShow("")
		}
	}

	return newYaos, cp
}

// 格式化输出
func (dt *TaoBstManager) FormatShow(cont string) string {
	/*
		:param cont:  需要显示的内容
		:return:
	*/
	spaceStr := []string{} //{" ", dt.Indent}
	for i := 0; i < dt.Indent; i++ {
		spaceStr = append(spaceStr, " ")
	}
	msg := strings.Join(spaceStr, "") + cont
	print(msg)
	return msg + "\n"
}

// # 6爻卦象 显示原始卦象, 返回卦值 和 第几卦
func (dt *TaoBstManager) KanGuaOrigin(env string) ([]int, int) {

	fmt.Printf("%v", start)
	guas, bstm := dt.YaoWithReBuild(49, env)
	fmt.Println("卦象已出:")

	yaos, postionChanges := dt.AnysisYaos(guas)

	fmt.Printf("yaos: postionChanges:, %v, %v, bstree:%v\n", yaos, postionChanges, bstm.Size)
	anysis, numb := ICApp.CommonText(yaos)
	fmt.Printf("第%v 卦, \n %v", numb, anysis)
	return yaos, numb

}

func DothTest(t int) {
	bstNew := &TaoBstManager{BinarySearchTree: &BinarySearchTree{}}
	if t <= 0 {
		t = 49
	}
	for i := 0; i < t; i++ {
		bstNew.Puts(i, fmt.Sprintf("suanzi_%v", i))
		// bstNew.Display()
		// time.Sleep(time.Second * 1)
		Logg.Printf("bstNew :%v, root:%v  len:%v\n  \n", bstNew, bstNew.Root, bstNew.Size)

	}

	// 以人为本 分天地人 查找某个节点 并以此旋转 再平衡
	rootKey := DoRand(t)
	keyNodes := bstNew.Searcher(rootKey)
	// Logg.Printf("keyNodes:%v\n", keyNodes)

	//根据bst2 再造新树
	newTree := bstNew.Rebuild(keyNodes)
	// oneTotal := newTree.Size
	caches := newTree.IterCache()
	Logg.Printf("newTree :%v, root:%v caches len:%v\ncaches:%v\n", newTree.Size, newTree.Root, caches.Size, len(caches.Read))

	// //根节点作为 人才
	Logg.Printf("人 newTree.Root:%v\n", newTree.Root)
	// panic("")
	// //查看左子树 天才
	// //查看右子树 地才

	// Logg.Printf("left size:%v, right size:%v\n", newTree.LeftSize, newTree.RightSize)

	Logg.Printf("before letf:%v   right:%v   \n", newTree.LeftSize, newTree.RightSize)

	//找到最小和最大的，以便预知随机生成后 树的形状
	mi := newTree.Root.FindMin() //最小
	mx := newTree.Root.FindMax() //最大
	st := newTree.Size           //总数
	//新的根节点 和以此生成的新的树 作为二变基础
	two := DoRange(mi.Key, mx.Key, "B")

	keyNodesTwo := newTree.Searcher(two)
	newTreeTwo := newTree.Rebuild(keyNodesTwo)

	//新天地 三变计算
	Logg.Printf("树大小应该没有变:%v, new:%v, equal:%v\n", st, newTreeTwo.Size, st == newTreeTwo.Size)

}
