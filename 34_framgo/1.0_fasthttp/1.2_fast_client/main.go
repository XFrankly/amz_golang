package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
 postwoman
 username=ma2002&password=man2012&price=12.13&total=1293

*/
// ClientConn.Client.Do(req)
type HttpClientHandler struct {
	Client http.Client
}

func ClientConnectionControl(AuthToken string, ContentType string) *HttpClientHandler {
	// 构造客户端 可 设置 请求头
	// 绑定 客户端， 调用头，参数
	var tokens string
	if AuthToken == "" {
		tokens = "" //fmt.Sprintf("Bearer %s", settings.AccessToken)
	} else {
		tokens = fmt.Sprintf("Bearer %s", AuthToken)
	}
	fmt.Printf("test tokens:%+v\n", tokens)
	if ContentType == "" {
		ContentType = "application/json"
	}
	hclients := &HttpClientHandler{
		Client: http.Client{},
	}

	return hclients
}

func main() {
	data := url.Values{}
	data.Set("username", "ma2002")
	data.Set("token", "123456789ma2002")
	data.Set("password", "man2012")
	data.Set("price", "12.22")
	data.Set("total", "1293")
	fmt.Printf("data:%+v\n", data)

	r, _ := http.NewRequest("POST", "http://localhost:8002/users/post", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	new_client := *ClientConnectionControl("", "x-www-form-urlencoded;charset=UTF-8")
	resp, err := new_client.Client.Do(r)
	if err != nil {
		msg := fmt.Sprintf("err request with r:%+v\n", r)
		panic(msg)
	}
	respBody, err0 := ioutil.ReadAll(resp.Body)
	if err0 != nil {
		fmt.Println("err to Read io", err0)
	}
	fmt.Printf("response:%+v \n, respBody.response:%+v\n", resp, string(respBody))
}
