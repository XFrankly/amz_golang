package main

/*
Hijacker 接口由 ResponseWriters 实现，它允许 HTTP 处理程序接管连接。

HTTP/1.x 连接的默认 ResponseWriter 支持 Hijacker，但 HTTP/2 连接故意不支持。
ResponseWriter 包装器也可能不支持 Hijacker。处理程序应始终在运行时测试此能力。
*/

import (
	"fmt"
	"log"
	"net/http"
)

/*
hijack
劫持讓調用者接管連接。
// 在調用劫持 HTTP 服務器庫之後 不會對連接做任何其他事情。
//
// 成為調用者的責任去管理並關閉連接。
//
// 返回的 net.Conn 可能有讀或寫截止日期 已經設置，具體取決於配置 服務器。
調用者有責任設置 或根據需要清除這些截止日期。

// 返回的 bufio.Reader 可能包含未處理的緩衝來自客戶端的數據。
//
// 調用 Hijack 後，原來的 Request.Body 一定不能使用。原始請求的上下文仍然有效並且
// 直到 Request 的 ServeHTTP 方法才被取消 返回。
*/
func main() {
	http.HandleFunc("/hijack", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()
		bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
		bufrw.Flush()
		s, err := bufrw.ReadString('\n')
		if err != nil {
			log.Printf("error reading string: %v", err)
			return
		}
		fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
		bufrw.Flush()
	})
}
