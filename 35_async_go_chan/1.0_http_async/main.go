package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

/*
使用chan 实现 golang 异步请求， 函数与 异步调用 有通信
1 将返回类型设置为通道。
2 创建一个通道以返回一个值。
3 为异步执行添加一个 go func。
4 在函数中，将值分配给通道。
5 在函数的最后，用值表示通道的返回。
6 在 main 函数中，将通道的返回值赋给一个变量。
*/
var (
	wg = sync.WaitGroup{}
)

type WebState struct {
	Url   string
	State bool
	Resp  string
}

func GetHttp(url string, c chan WebState) {

	resp, err := http.Get(url)
	if err != nil {
		c <- WebState{Resp: string(resp.Status),
			Url:   url,
			State: false}
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Fatalln(err)
	}
	c <- WebState{Resp: string(body),
		Url:   url,
		State: true}
	fmt.Println(url, "...Done!", len(c))

	fmt.Println("len c:", len(c))
	// wg.Done()
}

// func awaitTask() <-chan WebState {
func awaitTask() []WebState {
	fmt.Println("Starting Task...")
	urls := []string{"http://127.0.0.1:18083/bussiness/list",
		"http://127.0.0.1:18083/bussiness/list"}

	//异步
	c := make(chan WebState)
	for _, url := range urls {
		go GetHttp(url, c)

	}
	result := make([]WebState, len(urls))
	for i, _ := range result {
		result[i] = <-c
		if result[i].State {
			fmt.Println(result[i].Url, "is up.")
		} else {
			fmt.Println(result[i].Url, "is down !!")
		}
	}
	return result
}

func main() {
	value := awaitTask()
	fmt.Println("await value:", value)

}
