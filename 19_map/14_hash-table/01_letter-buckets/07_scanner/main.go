package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// An artificial input source. 人工输入源的处理， 包括空格 87个字符
	const input = "It is not the critic who counts; not the man who points out how the strong man stumbles"
	scanner := bufio.NewScanner(strings.NewReader(input))
	// Set the split function for the scanning operation. 设置扫描操作的 分割动作
	scanner.Split(bufio.ScanWords)
	fmt.Println(len(scanner.Text()))
	// Count the words.  计算字数
	for scanner.Scan() {
		fmt.Println("scanner text", len(scanner.Text()), scanner.Text())  //输出 扫描内容
	}
	if err := scanner.Err(); err != nil {  // 如果扫描操作 报错 则输出系统 标准输出 内容
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
}
