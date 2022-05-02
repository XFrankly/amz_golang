package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

var (
	sum = 0
)

func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	///// 三元表达式 如果 condition 为 true，返回 trueVal，否则返回 falseVal
	if condition {
		return trueVal
	}
	return falseVal
}

func zero(z *int) {
	*z = 0

}

func testFunc(MyF func(int) int) {
	fmt.Println(MyF(8))
}

// ///闭包函数
func returnFunc(x string) func() {

	return func() {
		sum += 1
		fmt.Println(sum)
	}
}

func returnFuncInt(step int) func(int) int {

	return func(step2 int) int {
		sum += step + step2
		fmt.Println(sum)
		return sum
	}
}

///// 从map随机选择一个键的值
func RandMap() {
	/// 展平 map 并从中选择
	m := map[string]int{
		"A": 1,
		"B": 2,
	}
	mk := make([]string, 0, len(m))
	for k := range m {
		mk = append(mk, k)
	}
	fmt.Printf("%+v\n", mk)
	//// rand.Intn 它返回 [0,n) 范围内的随机整数
	fmt.Printf(mk[rand.Intn(len(mk))]) /// 永远选择最后一个值

	//// MapKeys 随机选择 map的键值
	for i := 0; i < 5; i++ {
		r := rand.Intn(len(m)) /// 在 map 长度范围内选择一个整数
		fmt.Println("r == ", r)
		for k := range m {
			if r == 0 {
				fmt.Println(" k  value:", k)
			}
			r--
		}
		// panic("unreachable.")
		rmk := reflect.ValueOf(m).MapKeys()
		fmt.Println(rmk) //// 返回map 字典的 键，排序不固定
		fmt.Println(rmk[rand.Intn(len(rmk))].Interface())
	}
}
func main() {
	a, b := 12, 3
	max := If(a > b, a, b).(int) //// 三元表达式
	fmt.Printf("%T\n", max)
	fmt.Println(max)

	zero(&a) ///// 如果变量 被传入指针参数的函数，此变量值将被改变
	fmt.Println(a)

	test1 := func(x int) int {
		return x * -1
	} //(8)
	fmt.Println(test1(9)) // 调用一个计算表达式
	testFunc(test1)       /// 通过函数调用函数

	returnFunc("hello")()      /// 1
	x := returnFunc("goodbye") /// 1
	x()

	returnFuncInt(5)(5)   //// 调用函数并计算  10
	y := returnFuncInt(6) // 返回一个可调用对象
	y(7)                  // 6+7 + 5 + 5 + 2

	RandMap()
}
