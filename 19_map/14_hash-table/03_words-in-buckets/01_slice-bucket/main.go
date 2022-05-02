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
	// 为保存计数 创建分片 桶 大小12
	buckets := make([][]string, 12)

	//此处的代码已从录音中更新
	// 见下面的解释

	// 遍历文字
	for scanner.Scan() {
		// 文字赋值
		word := scanner.Text()
		// 获取一个桶
		n := hashBucket(word, 12)
		// 往桶中 添加文字 桶
		buckets[n] = append(buckets[n], word)
	}
	// 打印每个桶的 大小
	for i := 0; i < 12; i++ {
		fmt.Println("len buckets: ",i, " - ", len(buckets[i]))
	}
	// 打印其中一个字词组中的 字词
	fmt.Println(buckets[6])
	fmt.Println("buckets len:", len(buckets))
	fmt.Println("cap buckets:", cap(buckets))
}

func hashBucket(word string, buckets int) int {
	var sum int
	for _, v := range word {
		sum += int(v)
	}
	// 返回 每个字词的 桶， 均匀桶
	//return sum % buckets
	//  启用以下 代码

	// 分布不均的桶
	return len(word) % buckets
}

/*
UPDATED CODE
Up above, the code has been updated from the recording
I used to have this ...
		buckets = append(buckets, []string{})
... and changed it to this ...
		buckets[i] = []string{}

REASON:

This line of code ...
	buckets := make([][]string, 12)
... creates a slice with len and cap equal to 12. I can now access each of the twelve positions in the slice by index and assign values to them. If I "append" to this slice, like this ....
		buckets = append(buckets, []string{})
... I am adding another twelve positions to my slice; my len increases to 24 and my cap increases to 24. This is unnecessary. I can, instead, just direclty begin accessing the first twelve positions in my slice ... and that's why I changed the code to this ...
		buckets[i] = []string{}

EVEN MORE EXPLANATION

You don't even need this entire chunk of code ...

	for i := 0; i < 12; i++ {
		buckets[i] = []string{}
	}

... as this code ...

	buckets := make([][]string, 12)

... creates a slice holding a []string, but it doesn't yet have a len or cap, so later I use append which is how you add an item to a slice in a position that does not yet have an item (beyond its current len).

Thank you to Lee Trent for pointing this out!
*/
