package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	Logg = log.New(os.Stderr, "INFO -:", 18)
)

func checkErr(err error) {
	if err != nil {
		Logg.Fatal(err)
	}
}

//// udp socket
func udp_client() {
	if len(os.Args) != 2 {
		Logg.Println("Usage:socket client message.", os.Args)
		os.Exit(1)
	}

	msg := os.Args[1]
	// 该示例向本地网络机器上的回显服务发送一条小消息。该消息被回显。
	//使用该Dial函数，我们在 本地UDP 网络的端口 8889 上为  系统创建一个套接字。

	con, err := net.Dial("udp", "127.0.0.1:8889") //  "debian:7"
	checkErr(err)
	defer con.Close()
	/// 将消息写入套接字Write
	wrst, err := con.Write([]byte(msg))
	Logg.Println("net write result:", wrst, err)
	checkErr(err)
	/// make 使该函数创建一个字节切片。然后我们创建对该切片的响应
	reply := make([]byte, 1024)
	rrst, err2 := con.Read(reply)

	Logg.Println("net read result:", rrst, err2)
	checkErr(err2)
	Logg.Println("reply:", string(reply))
}

/// tcp socket tcp套接字
// 服務是一個有用的調試和測量工具。當天服務的報價只是發送一條短消息，不考慮輸入。
/// 连接本地 8800 服务并发送 5个hello 消息，然后发送一个 STOP\n 通知服务退出
func tcp_client() {
	con, err := net.Dial("udp4", "127.0.0.1:9900")
	checkErr(err)
	defer con.Close()
	Logg.Printf("con info:%+v\n", con)
	msg := "STOP\n"
	msg_list := []string{}

	for i := 0; i < 50; i++ {
		msg_list = append(msg_list, fmt.Sprintf("Hello.%v\n", i))
	}
	// msg_list = append(msg_list, msg)
	Logg.Println(msg, msg_list)
	for _, msg := range msg_list {
		time.Sleep(time.Second * 1)
		wrst, err2 := con.Write([]byte(msg))
		checkErr(err2)
		Logg.Println("write result:", wrst, err2)
		reply := make([]byte, 1024)

		rrst, err3 := con.Read(reply)
		checkErr(err3)
		Logg.Println("read result:", rrst, err3, string(reply))
	}

}
func main() {
	tcp_client()
}
