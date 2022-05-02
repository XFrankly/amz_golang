package main

import (
	"fmt"
	"unsafe"
)
import "reflect"

func main() {
	foo := "a"
	for i := 250; i <= 340; i++ {
		fmt.Println(i, " - ", string(i), " - ", []byte(string(i)))
		fmt.Println(reflect.TypeOf(i), reflect.TypeOf(foo), []byte(foo))
	}


	fmt.Println(foo)
	fmt.Printf("%T \n", foo)
	fmt.Println(foo)

	test()
}

func test(){
	// 类型强制转换

	var a int =10
	var b *int =&a
	var c *int64 = (*int64)(unsafe.Pointer(b))
	//var d float64 = (float64)(unsafe.Pointer(b))
	fmt.Println(*c)

	// 接口类型判断 返回是否成立的断言语法
	var e interface{} = 11
	var f  interface{} = 11.11
	t, isf := f.(float64)
	if isf{
		fmt.Println("float64",t, f, isf )
	}
	t2, ok:= e.(int)
	if ok{
		fmt.Println("int",t2, e, "ok", ok)
	}

	// 类型检查
	fmt.Println(reflect.TypeOf(c))
	fmt.Println(reflect.TypeOf(*c))

	//switch assert type
	var x interface{} = 122
	switch i := x.(type) {
	case nil:
		fmt.Println("x is nil", x, i)                // type of i is type of x (interface{})
	case int:
		fmt.Println("x is int", x, i)                            // type of i is int
	case float64:
		fmt.Println("x is float64", x, i)                        // type of i is float64
	case func(int) float64:
		fmt.Println("x is func(int)", x, i)                       // type of i is func(int) float64
	case bool, string:
		fmt.Println("x is bool", x, i)  // type of i is type of x (interface{})
	default:
		fmt.Println("x type donot know", x, i)     // type of i is type of x (interface{})
	}


}
/*
NOTE:
Some operating systems (Windows) might not print characters where i < 256

If you have this issue, you can use this code:

fmt.Println(i, " - ", string(i), " - ", []int32(string(i)))

UTF-8 is the text coding scheme used by Go.

UTF-8 works with 1 - 4 bytes.

A byte is 8 bits.

[]byte deals with bytes, that is, only 1 byte (8 bits) at a time.

[]int32 allows us to store the value of 4 bytes, that is, 4 bytes * 8 bits per byte = 32 bits.
*/
