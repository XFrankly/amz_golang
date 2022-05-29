package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

/*
這個文件創建了主包，它聲明了 main() 函數。該函數將使用導入的包創建一個 UDP 客戶端。
main() 函數在 arguments 變量中收集命令行參數並包括錯誤處理。
UDP 客戶端使用常規文件 I/O 函數與 UDP 服務器交互。當您向 UDP 服務器發送 STOP 命令時，客戶端將終止。這不是 UDP 協議的一部分，但在示例中用於為客戶端提供退出方式。
net.ResolveUDPAddr() 函數返回一個 UDP 端點地址。 UDP 端點的類型為 UDPAddr，包含 IP 和端口信息。
使用 net.DialUDP() 函數建立與 UDP 服務器的連接。
bufio.NewReader(os.Stdin) 和 ReadString() 用於讀取用戶輸入。
ReadFromUDP() 函數從服務器連接讀取數據包，如果遇到錯誤將返回。
*/

var (
	Logg = log.New(os.Stderr, "INFO -", 18)
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		Logg.Println("Please give the host:port string")
		return
	}

	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		Logg.Println(err)
		return
	}
	Logg.Println("The UDP server is:", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		Logg.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			Logg.Println("Exiting UDP clint.")
			return
		}
		if err != nil {
			Logg.Println(err)
			return
		}
		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			Logg.Println(err)
			return
		}
		Logg.Println("Reply:", string(buffer[0:n]))
	}

}
