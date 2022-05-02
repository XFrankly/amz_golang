package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func main() {
	//  获取白鲸小说第10章
	res, err := http.Get("http://www.gutenberg.org/files/2701/old/moby10b.txt")
	if err != nil { //报错打印
		log.Fatal(err)
	}

	//  扫描页面内容
	scanner := bufio.NewScanner(res.Body)
	// 等待完成后 关闭扫描体
	defer res.Body.Close()
	// 设置分割操作函数到 扫描操作者
	scanner.Split(bufio.ScanWords)
	// 为保存计数 创建分片 桶
	buckets := make([]int, 12)
	// 遍历所有文字
	for scanner.Scan() {
		// 把所有文字分成12个 部分 装入桶中
		n := hashBucket(scanner.Text(), 12)
		buckets[n]++
	}
	// 打印 桶 的内容
	fmt.Println(buckets)
}

func hashBucket(word string, buckets int) int {
	var sum int
	// 按 每个文字内容 的int 汇总 来返回桶
	for _, v := range word {
		sum += int(v)
		fmt.Println("hashBucket letter:", sum, string(sum) )

	}

	bb := sum % buckets
	fmt.Println("hashBucket  :",   word, buckets, reflect.TypeOf(bb), bb)
	//返回
	return bb
}
