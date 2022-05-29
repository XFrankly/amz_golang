package src

import (
	"conditions/settings"
	"net/http"
)

type AuthError struct {
	Message string
}

type RestError struct {
	Message string
}

type HttpClientHandler struct {
	Client *http.Client
	// 请求头
	Header *http.Header
}

type IsRest struct {
	HttpOnly bool
}
type Cookies struct {
	Version     int // 默认0
	Name        string
	Value       interface{}
	Port        int
	Domain      string
	Path        string // 默认 /
	Secure      bool
	Expires     interface{} // 到期
	Discard     bool        //  是否丢弃默认 true
	Comment     interface{}
	Comment_url interface{}
	Rest        IsRest // 是否 rest
	Rfc2109     bool
}

///////////////////////////参数类 结构体
type LoginType struct {
	User     string
	Password string
}

func InitLoginType(user string, password string) *LoginType {
	if user == "" || password == "" {
		return &LoginType{
			User:     settings.Name,
			Password: settings.Passwords,
		}
	}
	return &LoginType{
		User:     user,
		Password: password,
	}

}
