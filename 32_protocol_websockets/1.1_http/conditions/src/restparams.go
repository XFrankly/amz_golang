package src

import (
	"conditions/resp"
	"conditions/settings"
	"encoding/json"
	"fmt"
	"os"
)

// 构建请求头参数
type RequestsHeaders struct {
	Host        string
	Port        int
	Protocol    string
	AccessToken string
}

// 泛型判断
type RestStruct interface {
	LoginType | resp.RespMap // 并集声明 int64  float64 本质上你正在将 联合函数声明 移动到新的类型约束
	// 当您希望 将 参数约束为 int64 或 float64 时，可用使用此Number类型约束 而不是写 int64 | float64
}

func BuildRequestParam(host string, port int, protocol string) *RequestsHeaders {
	if host == "" || port == 0 || protocol == "" {
		host, port, protocol = settings.Host, settings.Port, settings.Protocol //"127.0.0.1", 18083, "http"
	}
	// http.Request 对象 构造请求 的参数，返回参数结构体
	if settings.Tls == true {
		protocol = "https"
	}
	requests := RequestsHeaderController(host, port, protocol)
	return requests
}

func BuildRequestParamBody(json string) []byte {
	// 构造请求体 参数
	if json == "" {
		json = "{\"user\": \"postgre\",\"password\": \"postgre.2022\"}"
	}

	var arr2 = []byte(json)
	return arr2
}

//泛型支持需要 go 1.18+
func DecodeStructStr(group interface{}) string {
	// 结构体 转 为 混合类型 dict json
	// 如:
	//group := resp.RespMap{
	//	Status:  1,
	//	Message:  map[string]string{
	//		"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ",
	//		"user":"postgre",
	//	},
	//}

	if group != nil {
		b, err := json.Marshal(group)
		if err != nil {
			fmt.Printf("error:", err)
		}
		fmt.Printf("%T %v", b, b)
		os.Stdout.Write(b)
		return fmt.Sprintf("%s", b)
	} else {
		return ""
	}

}
func LoginParamStr(username string, passwd string) string {
	// 返回 LoginType 实例信息的 字符串形式
	// 构造一个Login 实例
	group := InitLoginType(username, passwd)
	// 结构体转为混合 类型json 字符串
	//group := resp.RespMap{
	//	Status:  1,
	//	Message:  map[string]string{
	//		"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ",
	//		"user":"postgre",
	//	},
	//}

	return DecodeStructStr(group)
}

func LoginJson(user string, password string) string {
	if user == "" || password == "" {
		return fmt.Sprintf("{\"user\":\"%v\",\"password\": \"%v\"}", settings.Name, settings.Passwords)
	}
	// 返回 自定义json 字符串
	return fmt.Sprintf("{\"user\":\"%v\",\"password\": \"%v\"}", user, password)
}
func (self *RequestsHeaders) Login(json string) []byte {
	// 构造请求体 body参数
	// 构造restful http 地址

	if json == "" {
		json = fmt.Sprintf("{\"user\":\"%v\",\"password\": \"%v\"}", settings.Name, settings.Passwords)
	}

	var arr2 = []byte(json)
	return arr2
	//return fmt.Sprintf("%s://%s:%v%s", self.Protocol, self.Host, self.Port, path)
}

func (self *RequestsHeaders) Url(path string) string {
	// 构造restful http 地址
	return fmt.Sprintf("%s://%s:%v%s", self.Protocol, self.Host, self.Port, path)
}

func (self *RequestsHeaders) GetWelcome() {
	return
}
