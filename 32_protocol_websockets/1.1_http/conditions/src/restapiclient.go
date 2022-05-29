package src

import (
	"conditions/settings"
	"fmt"
	//"io/ioutil"
	"net/http"
	//"net/http/httputil"
	//"strings"
)

// 作为http请求参数体 RequestsConnection 的构造器
func RequestsHeaderController(host string, port int, protocol string) *RequestsHeaders {
	// 作为 RequestsHeaders 的构造器
	s := &RequestsHeaders{
		Host:        host,
		Port:        port,
		Protocol:    protocol,
		AccessToken: settings.AccessToken,
	}
	// 作为链式调用
	return s
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

	if ContentType == "" {
		ContentType = "application/json"
	}
	hclients := &HttpClientHandler{
		Client: &http.Client{},
		Header: &http.Header{
			"Host":          []string{settings.Host},
			"Content-Type":  []string{ContentType},
			"Authorization": []string{tokens},
		},
	}

	return hclients
}
