package tasks

import (
	"bytes"
	"conditions/src"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	logg   = log.New(os.Stderr, "[INFO] - ", 13)
	logger = log.New(os.Stderr, "[WARNING] - ", 13)
)

func httpRequest() {
	r := strings.NewReader("{\"id\": \"1\",\"title\": \"London\",\"artist\": \"PostBettyCarter\",\"price\": 79.99}") // post 的传入body的数据
	url := "http://127.0.0.1:18083/v1/booking/add"
	client := http.Client{}
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		//Handle Error
		logg.Println("http client err", err)
	}

	req.Header = http.Header{
		"Host":          []string{"www.host.com"},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer Token"},
	}

	res, err := client.Do(req)
	if err != nil {
		//Handle Error
		logg.Println("do err", err, res)
	}
	body0, err0 := ioutil.ReadAll(res.Body)
	logg.Printf("err:%s,body: %s \n", err0, body0)
}
func io_read_all_from_api_booking() {
	// 从自定义接口 读取并返回数据
	// Reader  作为数据源
	r := strings.NewReader("{\"id\": \"1\",\"title\": \"London\",\"artist\": \"PostBettyCarter\",\"price\": 79.99}") // post 的传入body的数据
	url := "http://127.0.0.1:18083/v1/booking/add"
	//client := http.Client{}
	//req, err := http.NewRequest("POST", url, r)
	//if err != nil {
	//	//Handle Error
	//	logg.Println("http client err", err)
	//}
	//
	//req.Header = http.Header{
	//	"Host":          []string{"www.host.com"},
	//	"Content-Type":  []string{"application/json"},
	//	"Authorization": []string{"Bearer Token"},
	//}
	//
	//res, err := client.Do(req)
	//if err != nil {
	//	//Handle Error
	//	logg.Println("do err", err, res)
	//}
	//body0, err0 := ioutil.ReadAll(res.Body)
	//logg.Printf("err:%s,body: %s \n", err0, body0)

	resp, err := http.Post(url, "applicaation/json", r)
	// req, err := http.NewRequest("POST", "http://127.0.0.1:8083/v1/booking", nil)
	// 关闭 响应包Body defer resp.Body.Close()
	body, err1 := io.ReadAll(resp.Body) // 必须io read
	logg.Printf("%T \n", body)
	body_str := fmt.Sprintf("%s", body)
	logg.Printf("POST resp type %T, resp.Body %T, %T, %s \n", resp, &resp.Body, body_str, body)
	logg.Println("POST /v1/booking resp: ", resp, "\nBody:", body_str, len(body), "\n", *resp, "\nerr:", err, err1)

	r2 := strings.NewReader("")
	respc, errcc := http.Post("http://127.0.0.1:18083/v1/booking/clean", "applicaation/json", r2)
	if errcc != nil {

		logg.Println("err post clean:", respc, errcc)
	}
	defer resp.Body.Close()
	//defer respc.Body.Close()
	bodyc, errc1 := io.ReadAll(resp.Body)   // 必须io read
	respcc, errc2 := io.ReadAll(respc.Body) // 必须io read
	bodyc_str := fmt.Sprintf("%s", bodyc)
	logg.Println(respcc, errc2)
	logg.Printf("%T bodyc_str:%s, Done. \n", bodyc, bodyc_str)
	logger.Printf("POST bodyc type %T, resp.Body %s,%v \n", bodyc, bodyc_str, errc1)
}

func ReturnResponseBody(r *http.Response) []byte {
	// 读取 *http.response 并返回 []byte
	body0, err0 := ioutil.ReadAll(r.Body)
	if err0 != nil {
		logg.Fatal("err readall with respone body: %s \n", err0)
	}

	logg.Printf("Response:type %T \n", body0)
	return body0
}
func ClientPipePath() {
	//io.pipe()  可以用于大量数据流的处理
	//pr, rw := op.Pipe() // 返回Reader 和Writer接口实现对象，利用这个特性就开头实现流式写入
	// 创建一个协程来写入，然后把Reader传递到方法，从而实现 http client body 的流式写入
	buf, rw := io.Pipe() // 会使用chunked 编码 进行传输 Transfer-Encoding: chunked
	// 提前计算出ContentLength并且对性能要求比较 苛刻时，手动设置 ContentLength优化性能
	// 开协程写入大量数据
	go func() {
		for i := 0; i < 10000; i++ {
			rw.Write([]byte(fmt.Sprintf("line:%d\r\n", i)))
		}
		rw.Close()
	}()
	// 传递Reader
	http.Post("http://127.0.0.1:18083/v1/booking", "text/pain", buf)
}
func ClientRequestPath(requests *src.RequestsHeaders, method string, path string, body []byte) *http.Response { //
	// 组合 url 和 参数，发起请求
	hClients := src.ClientConnectionControl("", "")
	fullPath := requests.Url(path)

	// 构建客户端
	if body == nil {
		// 使用默认的 用户名密码登录
		body = requests.Login("")
	}
	//io.pipe()  可以用于大量数据流的处理
	//pr, rw := op.Pipe() // 返回Reader 和Writer接口实现对象，利用这个特性就开头实现流式写入
	// 创建一个协程来写入，然后把Reader传递到方法，从而实现 http client body 的流式写入
	bytesBody := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, fullPath, bytesBody) //bytes.NewBuffer(body))

	logg.Println(method, "<-- URL ->> :", fullPath, "request body:", bytesBody)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", fullPath, "err:", err, method, fullPath)
	}

	req.Header.Set("X-Custom-Header", "FuckTheWorld.")
	req.Header.Set("Content-Type", "application/json")
	req.Header = *hClients.Header
	hClients.Client.Timeout = time.Second * 250 //超时时间 25s
	res, err := hClients.Client.Do(req)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", fullPath, "err:", err, method, fullPath)
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

func ClientRequestPathByReader(requests *src.RequestsHeaders, method string, path string, body *strings.Reader) {
	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	hclients := src.ClientConnectionControl("", "")
	full_path := requests.Url(path)
	req, err := http.NewRequest(method, full_path, body) //bytes.NewBuffer(body))
	logger.Println(method, "<-- URL ->> :", full_path, "ClientRequestPath err", err, method, full_path)
	if err != nil {
		//Handle Error
		logger.Println("ClientRequestPath err", err, method, full_path)
	}

	req.Header.Set("X-Custom-Header", "FuckTheWorld.")
	req.Header.Set("Content-Type", "application/json")
	req.Header = *hclients.Header
	hclients.Client.Timeout = time.Second * 250 //超时时间 25s
	res, err := hclients.Client.Do(req)
	if err != nil {
		//Handle Error
		logger.Println("do err", err, res)
	}
	body0, err0 := ioutil.ReadAll(res.Body)
	logg.Printf("err:%s,body: %s \n", err0, body0, method, full_path)
	logg.Println(full_path)
}
