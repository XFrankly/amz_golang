package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
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
	buckets := make([]int, 200)
	// 遍历所有文字
	for scanner.Scan() {
		n := hashBucket(scanner.Text())
		buckets[n]++
	}
	// 截取 桶 切片的部分内容， 内容 比 原文字小
	fmt.Println(buckets[65:123])

	// 遍历并查看 桶中的内容
	fmt.Println("***************")
	for i := 28; i <= 126; i++ {
		fmt.Printf("%v - %c - %v \n", i, i, buckets[i])
		fmt.Println(string(buckets[i])) //查看 字符串内容

	}
}

func hashBucket(word string) int {
	return int(word[0])
}
