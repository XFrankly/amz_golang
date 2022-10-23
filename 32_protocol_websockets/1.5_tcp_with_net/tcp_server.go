package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	Logg = log.New(os.Stderr, "INFO -", 18)
)

/*
這個文件創建了主包，它聲明了 main() 函數。該函數將使用導入的包來創建 TCP 服務器。
main() 函數在 arguments 變量中收集命令行參數並包括錯誤處理。
net.Listen() 函數使程序成為 TCP 服務器。此函數返回一個 Listener 變量，它是面向流協議的通用網絡偵聽器。
只有在成功調用 Accept() 之後，TCP 服務器才能開始與 TCP 客戶端交互。
TCP 服務器的當前實現只能為連接到它的第一個 TCP 客戶端提供服務，因為 Accept() 調用在 for 循環之外。
在本指南的創建並發 TCP 服務器部分，您將看到一個 TCP 服務器實現，它可以使用 Goroutine 為多個 TCP 客戶端提供服務。
TCP 服務器使用常規文件 I/O 函數與 TCP 客戶端交互。這種交互發生在 for 循環內。
與 TCP 客戶端類似，當 TCP 服務器收到 TCP 客戶端的 STOP 命令時，它會終止。

... 启动服务 go run .\tcp_server.go 8899
将在接收到客户端 的消息后 显示客户端发送是消息
*/
func main() {
	// arguments := os.Args
	// if len(arguments) == 1 {
	// 	Logg.Println("Please provide port number")
	// 	// return
	// 	arguments[1] = "8800"
	// }

	PORT := ":" + "8800" //arguments[1]
	listens, err := net.Listen("tcp", PORT)

	Logg.Printf("tcp listen port %+v listens:%+v\n", PORT, listens)
	if err != nil {
		Logg.Println(err)
		return
	}
	defer listens.Close()

	c, err := listens.Accept()
	Logg.Printf("accept:%+v, err:%+v\n", c, err)
	if err != nil {
		Logg.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		Logg.Println("net data reader:", netData, c, err)

		if err != nil {
			Logg.Println(err)
			return
		}

		receive_msg := strings.TrimSpace(string(netData))
		Logg.Printf("tcp server receive_msg:%+v\n", receive_msg)
		if receive_msg == "STOP" {
			Logg.Println("Exiting TCP server!")
			return
		}

		Logg.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}
}
