package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 该地址是 开源小说站点，白鲸（Moby-Dick） 是赫尔曼·梅尔维尔（Herman Melville）的小说， 于1851年首次出版
	res, err := http.Get("http://www.gutenberg.org/files/2701/old/moby10b.txt")
	// get获取url 内容
	if err != nil {   // 如果有报错
		log.Fatal(err)  //打印错误
	}
	bs, err := ioutil.ReadAll(res.Body) // io读取 get 获取的 body内容
	res.Body.Close()   // 获取内容对象关闭
	if err != nil { // 获取内容时报错
		log.Fatal(err)  //打印错误日志
	}
	fmt.Printf("%s", bs)  //输出 获取的内容
}
