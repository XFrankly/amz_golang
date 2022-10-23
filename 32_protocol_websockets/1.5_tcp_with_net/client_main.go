package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/*
這個文件創建了主包，它聲明了 main() 函數。該函數將使用導入的包創建 TCP 客戶端。
main() 函數在 arguments 變量中收集命令行參數，並確保發送了 host:port 的值。
CONNECT 變量存儲要在 net.Dial() 調用中使用的 arguments[1] 的值。
對 net.Dial() 的調用開始 TCP 客戶端的實現，並將您連接到所需的 TCP 服務器。
net.Dial() 的第二個參數有兩個部分；第一個是 TCP 服務器的主機名或 IP 地址，第二個是 TCP 服務器偵聽的端口號。
bufio.NewReader(os.Stdin) 和 ReadString() 用於讀取用戶輸入。
任何用戶輸入都使用 Fprintf() 通過網絡發送到 TCP 服務器。
bufio reader 和 bufio.NewReader(c).ReadString('\n') 語句讀取 TCP 服務器的響應。
為簡單起見，此處忽略錯誤變量。
用於讀取用戶輸入的整個 for 循環只會在您向 TCP 服務器發送 STOP 命令時終止。

... 连接服务器 go run .\client_main.go 127.0.0.1:8899
>> HELLO
>> STOP  // 这将退出tcp 服务器 和 客户端
*/

func conn_tcp_server() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		///// 這些例程以 'f' 結尾並採用格式字符串。
		// Fprintf 根據格式說明符格式化並寫入 w.
		// 它返回寫入的字節數和遇到的任何寫入錯誤。
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

func main() {
	conn_tcp_server()
}
