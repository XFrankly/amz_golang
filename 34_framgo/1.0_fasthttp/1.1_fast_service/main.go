package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

/*

 */
type RequestContext struct {
	user     string
	password string
}

func httpHandle(ctx *fasthttp.RequestCtx) {
	//获取请求体的body
	body := ctx.Request.Body()

	//输出
	fmt.Println(string(body))

	// 其他业务处理,接收参数body
	reqContext := &RequestContext{}
	err := json.Unmarshal(body, reqContext)
	if err == nil {
		fmt.Println("err nil")
	}
	ctx.Response.AppendBodyString("ok")
	ctx.Response.SetStatusCode(200)
	fmt.Printf("ctx: %+v\n", ctx)
	fmt.Printf("reqContext:%+v body:%+v\n", *reqContext, string(body))

}
func main() {
	port := 8002
	fmt.Printf("server listen:%d\n", port)
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		switch path {
		case "/users/post":
			httpHandle(ctx)
		default:
			fmt.Println("==============================")
		}
	}))
}
