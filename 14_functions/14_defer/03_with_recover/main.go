package main

import "fmt"

func main() {
	f1()
	f2()
	f3()
}

func f1() {
	fmt.Println("func f1")
}

//捕获panic
func f2() {
	// defer 一定要在 panic 之前, 因为 panic 触发时
	// panic 所在函数就会立刻终止并倒序调用所有已经存在的defer
	// 若 defer 在 panic 之后, 程序根本不会知道后面有 defer
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err =", err)
			fmt.Println("recover in f2 (first)")
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err =", err)
			fmt.Println("recover in f2 (second)")
		}
	}()
	panic("panic in f2")
}

func f3() {
	fmt.Println("func f3")
	panic("panic func f3")
}
