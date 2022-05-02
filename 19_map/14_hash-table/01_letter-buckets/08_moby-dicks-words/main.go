package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// get the book moby dick 获取 小说 白鲸 文本
	res, err := http.Get("http://www.gutenberg.org/files/2701/old/moby10b.txt")
	if err != nil {  // 报错则输出 错误内容
		log.Fatal(err)
	}

	// scan the page 扫描查询该页
	// NewScanner takes a reader and res.Body implements the reader interface (so it is a reader)
	// NewScanner接收阅读器并执行res.Body实现阅读器接口（因此它是阅读器）
	scanner := bufio.NewScanner(res.Body)
	// 等待扫描完成后 关闭 读取接口
	defer res.Body.Close()
	// Set the split function for the scanning operation.
	// 设置扫描操作的分割功能。
	scanner.Split(bufio.ScanWords)
	// Loop over the words 遍历单词
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
