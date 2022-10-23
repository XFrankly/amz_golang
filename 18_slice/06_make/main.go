package main

import (
	"fmt"
	"unsafe"
)

//通过指针获取切片某个位置的值
func SliceValue() {
	// 根据内存地址获取下一个字节内存地址对应的值
	dataList := [3]int8{11, 22, 33}

	// 1. 获取数组第一个元素的地址
	var firstDataPtr *int8 = &dataList[0]

	// 2. 转换成Pointer类型
	ptr := unsafe.Pointer(firstDataPtr)

	// 3. 转换成uIntPtr类型，然后进行内存地址的计算（即：地址加一个字节，意味着取第二个索引位置的值）
	targetAddress := uintptr(ptr) + 2

	// 4. 根据新地址，重新转换成Pointer类型
	newPtr := unsafe.Pointer(targetAddress)

	// 5. Pointer对象转换为int8指针类型
	value := (*int8)(newPtr)

	// 6. 根据指针获取值
	fmt.Println("最终结果为：", *value)
}

func CompareSlice(s1 []string, s2 []string) bool {
	if s1 == nil && s2 == nil {
		return true
	} else if s1 == nil || s2 == nil {
		return false
	}

	// for i, s := range s1 {
	// 	for j, k := range s2 {
	// 		if s == k {
	// 			s1 = s1[:i] + s1[i+1:]
	// 		}
	// 	}
	// }
	return false
}
func main() {

	customerNumber := make([]int, 3)
	// 3 is length & capacity
	// // length - number of elements referred to by the slice
	// // capacity - number of elements in the underlying array
	customerNumber[0] = 7
	customerNumber[1] = 10
	customerNumber[2] = 15

	fmt.Println(customerNumber[0])
	fmt.Println(customerNumber[1])
	fmt.Println(customerNumber[2])

	greeting := make([]string, 3, 3)
	// 3 is length - number of elements referred to by the slice 切片所引用的元素数
	// 5 is capacity - number of elements in the underlying array 基础数组中的元素数
	// you could also do it like this

	greeting[0] = "Good morning!"
	greeting[1] = "Bonjour!"
	greeting[2] = "dias!"

	fmt.Println(greeting[2])
	fmt.Println(greeting, len(greeting))
	fmt.Println(greeting[1:]) //保留切片中 1及其以后的
	fmt.Println(greeting[0:]) //保留切片中 0及其以后的
	newG := greeting[:]       //保留切片全部
	fmt.Println(newG)
	fmt.Println(greeting[:1]) //保留切片中的前一个
	fmt.Println(greeting[:2]) //保留切片中的 前两个

	SliceValue()
}
