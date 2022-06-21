package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func httpHandle(ctx *fasthttp.RequestCtx) {
	//获取请求体的body
	body := ctx.Request.Body()

	//输出
	fmt.Println(string(body))

	// 其他业务处理,接收参数body
	reqContext := &model.RequestContext{}
	err := json.Unmarshal(body, reqContext)
	if err == nil {
		fmt.Println("err nil")
	}
	ctx.Response.AppendBodyString("ok")
	ctx.Response.SetStatusCode(204)
}
func main() {
	log.Fatal(fasthttp.ListenAndServe(":8002", func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		switch path {
		case "/post":
			httpHandle(ctx)
		default:
			fmt.Println("==============================")
		}
	}))
}
