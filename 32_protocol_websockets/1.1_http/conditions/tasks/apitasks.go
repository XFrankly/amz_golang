package tasks

import (
	"bytes"
	"conditions/resp"
	"conditions/settings"
	"conditions/src"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 定义 http 任务
type Tasks struct {
	ClientParam *src.RequestsHeaders
	ClientConn  *src.HttpClientHandler
}

//func (self *Tasks) DefaultAuthTask(loginjson string) []byte {
func DefaultAuthTask(loginjson string) []byte {
	if loginjson == "" {
		loginjson = fmt.Sprintf("{\"user\":\"%v\",\"password\": \"%v\"}", settings.Name, settings.Passwords)
	}
	//"http://127.0.0.1:18083/v1/booking/clean"   http://localhost:18083/v1/loginjson
	//book_path := "/v1/booking/add"
	loginpath := "/v1/loginjson"
	host, port, protocol, method := "127.0.0.1", 18083, "http", "POST"
	requestsParam := src.BuildRequestParam(host, port, protocol)

	//var userPasswd = src.BuildRequestParamBody(strJson)
	userPasswd := requestsParam.Login(loginjson)
	httpResponse := ClientRequestPath(requestsParam, method, loginpath, userPasswd)

	return ReturnResponseBody(httpResponse) // 将在c.Keys里，不在body
}

func (self *Tasks) AuthTask(loginjson string) []byte {
	if loginjson == "" {
		loginjson = fmt.Sprintf("{\"user\":\"%v\",\"password\": \"%v\"}", settings.Name, settings.Passwords)
	}
	return DefaultAuthTask(loginjson)
}

func GetAuthToken(loginjson string) string {
	if loginjson == "" {
		loginjson = src.LoginParamStr("", "")
	}

	loginResp := DefaultAuthTask(loginjson)

	respMap := resp.InitRespMap()
	objresp := resp.ReStrToStruct(loginResp, respMap)

	logger.Println("Done.", objresp.Status, objresp.Message["token"], objresp.Message["user"])
	return objresp.Message["token"]
}

func (self *Tasks) TaskLoginParamStr() string {
	return src.LoginParamStr("", "")
}
func (self *Tasks) AuthToken(loginjson string) string {
	if loginjson == "" {
		loginjson = self.TaskLoginParamStr()
	}

	loginResp := self.AuthTask(loginjson)

	respMap := resp.InitRespMap()
	objresp := resp.ReStrToStruct(loginResp, respMap)
	logger.Println("\n", objresp.Status, objresp.Message["token"], objresp.Message["user"], "\n")
	logger.Println("Auth Done.")
	return objresp.Message["token"]
}

func (self *Tasks) ConnectionControl(AuthToken string) {
	ClientHandler := src.ClientConnectionControl(AuthToken, "")
	self.ClientConn = ClientHandler
}

func (self *Tasks) Do_Request(method string, full_path string, body []byte) *http.Response {
	bytes_body := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, full_path, bytes_body) //bytes.NewBuffer(body))

	logg.Println(method, "<-- request ->> :", "request body:", bytes_body)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", full_path, "err:", err, method, full_path)
	}

	req.Header.Set("X-Custom-Header", "FuckTheWorld.")
	req.Header.Set("Content-Type", "application/json")
	req.Header = *self.ClientConn.Header

	self.ClientConn.Client.Timeout = time.Second * 250 //超时时间 25s
	res, err := self.ClientConn.Client.Do(req)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", full_path, "\n err:", err, method)
	}
	// res 属于 *http.Response 结构如下
	//Status     string // e.g. "200 OK"
	//	StatusCode int    // e.g. 200
	//	Proto      string // e.g. "HTTP/1.0"
	//	ProtoMajor int    // e.g. 1
	//	ProtoMinor int    // e.g. 0
	//Header Header
	//Body io.ReadCloser
	//ContentLength int64
	//TransferEncoding []string
	//Close bool
	//Uncompressed bool
	//Trailer Header
	//Request *Request
	//TLS *tls.ConnectionState

	return res

}
func (self *Tasks) BookingPostPageClean(makeAuth bool) *http.Response {
	// /booking/:action
	postPath := "/v1/booking/clean"
	var AuthToken string
	// 获取 token
	if makeAuth == true {
		AuthToken = self.AuthToken("")
	}

	fullUrl := self.ClientParam.Url(postPath)
	// 绑定客户端
	self.ConnectionControl(AuthToken)
	// 执行调用
	var eb = []byte(`{}`)

	httpResponse := self.Do_Request("POST", fullUrl, eb)
	body0, err0 := ioutil.ReadAll(httpResponse.Body)
	if err0 != nil {
		logger.Println("err to Read io", err0)
	}
	logg.Printf("%s \n", fmt.Sprintf("%s", body0))
	//logg.Println(httpResponse)
	return httpResponse // 将在c.Keys里，不在body
}

func (self *Tasks) BookingPostPage(makeAuth bool) *http.Response {
	// /booking/:action
	postPath := "/v1/booking/add"
	var AuthToken string
	// 获取 token
	if makeAuth == true {
		AuthToken = self.AuthToken("")
	}

	fullUrl := self.ClientParam.Url(postPath)
	// 绑定客户端
	self.ConnectionControl(AuthToken)
	// 执行调用
	var eb = []byte(`{"id": "5","title": "London","artist": "BettyCarter","price": 149.99}`)

	httpResponse := self.Do_Request("POST", fullUrl, eb)
	body0, err0 := ioutil.ReadAll(httpResponse.Body)
	if err0 != nil {
		logger.Println("err to Read io", err0)
	}
	logg.Printf("%s \n", fmt.Sprintf("%s", body0))
	//logg.Println(httpResponse)
	return httpResponse // 将在c.Keys里，不在body
}

func (self *Tasks) Closer() {
	self.ClientConn.Header.Set("Authorization", "")
	self.ClientParam.AccessToken = ""
}

// 调用 入口
func DoTasks() *Tasks {
	return &Tasks{
		ClientParam: src.BuildRequestParam("", 0, ""),
		ClientConn:  src.ClientConnectionControl("", ""),
	}
}

//
//
//func ClientRequestPath(host string, port int, protocol string, method string, path string) {
//	r := strings.NewReader("{\"id\": \"1\",\"title\": \"London\",\"artist\": \"PostBettyCarter\",\"price\": 79.99}") // post 的传入body的数据
//
//	requestsParam, hclients := src.ClientConnectionControl(host, port, protocol)
//	full_path := requestsParam.Url(path)
//	req, err := http.NewRequest(method, full_path, r)
//	if err != nil {
//		//Handle Error
//		logger.Println("ClientRequestPath err", err, method, full_path)
//	}
//	req.Header = *hclients.Header
//	res, err := hclients.Client.Do(req)
//	if err != nil {
//		//Handle Error
//		logger.Println("do err", err, res)
//	}
//	body0, err0 := ioutil.ReadAll(res.Body)
//	logg.Printf("err:%s,body: %s \n", err0, body0, method, full_path)
//	logg.Println(full_path)
//}
