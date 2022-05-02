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
	// 取第一个字
	letter := int(word[0])
	fmt.Println("hashBucket letter:", letter, string(letter), word, buckets)

	bucket := letter % buckets
	fmt.Println("hashBucket:", bucket, string(bucket), bucket, buckets)

	// 返回
	return bucket
}
