package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 获取夏洛克·福尔摩斯的冒险书
	res, err := http.Get("http://www.gutenberg.org/cache/epub/1661/pg1661.txt")
	if err != nil { //打印错误
		log.Fatal(err)
	}

	// 扫描页面
	scanner := bufio.NewScanner(res.Body)
	// 完成后关闭扫描 体
	defer res.Body.Close()
	// 设置分割操作函数到 扫描操作者
	scanner.Split(bufio.ScanWords)
	// 以int为键创建 哈希字典
	// 和另一张 哈希字典的值
	// 代表字符串的键，这将是单词
	// 和一个int值，这将是单词出现的次数
	buckets := make(map[int]map[string]int)
	// 创建切片以容纳单词
	for i := 0; i < 12; i++ {
		buckets[i] = make(map[string]int)
	}
	// 遍历扫描体
	for scanner.Scan() {
		// 赋值给 word
		word := scanner.Text()
		// 获取hash桶
		n := hashBucket(word, 12)
		// 更新 hash 键值
		buckets[n][word]++
	}
	// 打印桶中 单词
	for k, v := range buckets[6] {
		// 出现次数， 单词本身
		fmt.Println(v, " \t- ", k)
	}
}

func hashBucket(word string, buckets int) int {
	var sum int
	for _, v := range word {
		sum += int(v)
	}
	// 均匀桶
	return sum % buckets
}
