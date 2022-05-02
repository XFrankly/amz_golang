package main

import "fmt"

// 声明结构体 a struct
type Comic struct {
	// 声明 结构体属性
	Universe string
}

// 基础结构体的函数
func (comic Comic) ComicUniverse() string {
	// 返回 comic 结构体
	return comic.Universe
}

// 另一个结构体声明
type Marvel struct {
	// 匿名区， 结构体 内嵌其他结构体的地方
	Comic
}
type DC struct {
	// 匿名区， 结构体 内嵌其他结构体的地方
	Comic
}

func DoFirst() {
	// 创建实例
	c1 := Marvel{
		// 子结构体 可以访问 基础结构体 属性
		Comic{Universe: "MCU"},
	}
	fmt.Println("Universe is:", c1.ComicUniverse())

	c2 := DC{
		Comic{Universe: "DC"},
	}
	fmt.Println("New Universe is:", c2.ComicUniverse())
}

////////////////////////////// case 2
type first struct {
	// 声明 基本结构体1 属性
	base_one string
}

type second struct {
	//声明 基本结构体2 属性
	base_two string
}

// 基本结构体1 函数 返回结构体属性
func (f first) printBase1() string {
	// 返回结构体 属性字符串
	return f.base_one
}

func (s second) printBase2() string {
	// 返回结构体2 属性 字符串
	return s.base_two
}

// 子结构体 内联多个 基础结构体
type child struct {
	// 匿名区， 内联多个基础结构体
	first
	second
}

// 子结构体调用
func DoSeconds() {
	cc := child{
		first{base_one: "In base struct 1."},
		second{base_two: "\nIn base struct 2. \n"},
	}
	fmt.Println(cc.printBase1())
	fmt.Println(cc.printBase2())
}

/////////////// 接口 类型 结构体的 继承, 子类型化
type iBase interface {
	say()
}

type base struct {
	color string
	/// 解决 子类型 和 基础类型 有相同函数 名的问题
	clear func()
	news  interface{}
}

func (b *base) say() { /// 函数是 Golang的一等变量
	// fmt.Println("Hi from say function.")
	b.clear()
}

// /// 解决继承问题，基础类型 已经包含了一个 函数属性
// func (b *base) clear() { // 继承类型的 清理函数
// 	fmt.Println("Clear from base's function.")
// }

type childthree struct {
	base  // embedding
	style string
}

func (c *childthree) clear() {
	fmt.Println("Clear from childthree's function.")
}

func check(b iBase) { // 这里可以接收 child 实例
	b.say()
}

func DoThreed() {
	base0 := base{color: "Blue", clear: func() { fmt.Println("clear sub child base0 func") }}
	base := base{color: "Red",
		//  实例化时 定义一个 函数 给 子类型
		clear: func() {
			fmt.Println("Clear from child's base func")
		},
		news: base0,
	}

	child3 := &childthree{
		base:  base,
		style: "somestyle",
	}
	child3.say() //  child 类型的 say 调用了 基础类的 clear
	fmt.Println("The color is "+child3.color, child3.news)
	check(child3) // iBase 基础类型的 say
}

////// 使用 interface 实现 多重继承
// type iBase1 interface {
// 	say()
// }
// type iBase2 interface {
// 	walk()
// }

type base1 struct {
	color string
}

func (b *base1) say() {
	fmt.Println("Hi from say function.", b.color)
}

func (b *base1) walk() {
	fmt.Println("Hi from walk function.", b.color)
}

type base2 struct {
	language string
}

func (b *base2) clear() {
	fmt.Println("Clear from base2's function.", b.language)
}

func (b *base2) talk() {
	fmt.Println("Hi from talk language:", b.language)
	// b.clear()

}

type childFour struct {
	base1 // 内联
	base2 // 内联
	style string
	// clear func()
}

func (b *childFour) clear() {
	fmt.Println("Clear from childFour's function.", b.style, b.language)
}

// func checkSay(b iBase1) {
// 	b.say()
// }
// func checkWalk(b iBase2) {
// 	b.walk()

// }

func DoFourMulinherit() {
	base1 := base1{color: "Red"}
	base2 := base2{language: "En, Sp, Jp"}
	child4 := &childFour{
		base1: base1,
		base2: base2,
		style: "Child,somestyle",
	}
	child4.say()
	child4.walk()
	child4.talk()
	child4.clear()
	fmt.Printf("child info:%+v\n", child4)
	fmt.Println("child4 attr:", child4.color, child4.language, child4.style)
	// checkSay(child4)
	// checkWalk(child4)
}

////////////// 类型层次结构
type iAnimal interface {
	breathe()
}

type animal struct {
}

func (a *animal) breathe() {
	fmt.Println("Animal breate.")
}

type iAquatic interface {
	iAnimal
	swim()
}

type aquatic struct {
	animal
}

func (a *aquatic) swim() {
	fmt.Println("Aquatic swim.")
}

type iNonAquatic interface {
	iAnimal
	walk()
}
type nonAquatic struct {
	animal
}

func (a *nonAquatic) walk() {
	fmt.Println("Non-Aquatic walk.")
}

type shark struct {
	aquatic
}
type lion struct {
	nonAquatic
}

func checkAquatic(a iAquatic)       { a.swim() }
func checkNonAquatic(a iNonAquatic) { a.walk() }
func checkAnimal(a iAnimal)         { a.breathe() }
func DoFive() {
	shark := &shark{}
	checkAquatic(shark)
	checkAnimal(shark)
	lion := &lion{}
	checkNonAquatic(lion)
	checkAnimal(lion)

}

///////////////////////
func DoNothing() {
	////  return 将没有任何东西， nil 也不会返回
	return
}
func main() {
	DoFirst()
	fmt.Println("Two.+++++++++++++++++++++++++++++++")
	DoSeconds()
	fmt.Println("Three..###################......")

	DoThreed()
	fmt.Println("Four.........................")
	DoFourMulinherit()

	fmt.Println("Five@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	DoFive()

	fmt.Println("==================================================:")
	DoNothing()
	fmt.Println("five:")

	m1 := map[string]string{
		"A": "a1",
		"B": "a2",
		"C": "c3",
	}
	fmt.Printf("%+v\n", m1)
	fmt.Println("m1:", m1["A"], len(m1))
}
