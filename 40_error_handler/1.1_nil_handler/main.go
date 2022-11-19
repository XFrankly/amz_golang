package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type VideosService struct{}

/*
Nil在 Golang 中代表零值,很多 类型的值可能是nil，如 map [], 类型的指针等
硬编码nil总是表现得像未分配的接口，即 - (type=nil, value=nil)。这可以转化为另一种“奇怪”的行为，但是一旦您理解了它背后的机制，您将能够有效地使用它。
例子：
var p *int              // (type=*int,value=nil)
var i interface{}       // (type=nil,value=nil)

if i != nil {           // (type=nil,value=nil) != (type=nil,value=nil)

	    fmt.Println("not a nil")
	}

i = p                   // assign p to i

// a hardcoded nil is always nil,nil (type,value)
if i != nil {           // (type=*int,value=nil) != (type=nil,value=nil)

	    fmt.Println("not a nil")
	}

//the code outputs "not a nil" only once
最令人困惑的，可能会导致意外的程序行为，因为i变量可以传递给另一个函数，该函数将其interface{}作为输入类型参数，因此仅检查基本类型是i == nil不够的。解决办法是什么？

有两种解决方案，一种是将值与键入的 nil 值进行比较，第二种是使用反射包。看看这个例子

该表达式(*DefaultHandler)(nil)是从无类型的 nil 到 type 的零值的转换*DefaultHandler。
转换的形式是T(x)whereT是一个类型，x是一个可以转换为 type 的值T。在这个例子中，T是*DefaultHandler和x是nil。
需要括号*DefaultHandler来区分到指针类型的转换和对非指针类型的转换的取消引用。

表达式(*DefaultHandler)()不是有效的 Go 语法。
如果是具有复合文字语法的类型，&DefaultHandler{}也可以使用该值。DefaultHandler转换模式适用于所有类型。
该表达式 类似 &SomeType{}

	通常&SomeType{}会分配并初始化一个新值（除非编译器检测并优化它）
	(*SomeType)(nil)永远不会进行任何分配
*/

// 显式转换 等效于 &VideosService{}
var vid Handler = (*VideosService)(nil)

// 隐式转换
var _ Handler = (*VideosService)(nil)

type Handler interface {
}

// 对比空间分配占比
func CompareSpaceSize() {

	vidAddr := &VideosService{}
	fmt.Printf("size hander:%v\n", unsafe.Sizeof(vid))
	fmt.Printf("type of hander:%v\n", reflect.TypeOf(vid).String())
	fmt.Printf("type of nil in this module:%v\n", reflect.TypeOf(nil))
	fmt.Printf("size &hander:%v\n", unsafe.Sizeof(vidAddr))
}
func main() {
	//nil判断
	var i map[string]int
	var p *int
	var ss *string
	fmt.Printf("equal int nil?:%v, \nequal int point nil:%v,\nequal str point nil:%v\n",
		i == nil, p == nil, ss == nil)
	var ifs interface{}
	fmt.Printf("equal interface nil:%v\n", ifs == nil)

	// 不要直接判断 接口中的nil 是否nil，即使空接口中 添加了nil 值，那就不是空接口了
	var pp *int        // (type=*int,value=nil)
	var ii interface{} // (type=nil,value=nil)

	if ii != nil { // (type=nil,value=nil) != (type=nil,value=nil)
		fmt.Println("ii not a nil one")
	} else {
		fmt.Println("ii is a nil one")
	}

	ii = pp // assign p to i

	// a hardcoded nil is always nil,nil (type,value)
	if ii != nil { // (type=*int,value=nil) != (type=nil,value=nil)
		fmt.Println("ii not a nil two")
	} else {
		fmt.Println("ii is a nil two")
	}

	//初始化的空间使用对比
	CompareSpaceSize()
}
