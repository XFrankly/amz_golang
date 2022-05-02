package main

import "fmt"
import "reflect"

func main() {
	// 创建一个 包含1个元素的 分片
	buckets := make([]int, 1)
	fmt.Println("%T ", buckets)
	// 查看类型
	fmt.Println("reflect type", reflect.TypeOf(buckets), buckets)
	// 查看第0个元素的值
	fmt.Println(buckets[0])
	// 设置第0个元素的值
	buckets[0] = 42
	fmt.Println(buckets[0])
	// 改变第0个元素的值
	buckets[0]++
	fmt.Println(buckets[0])
}
