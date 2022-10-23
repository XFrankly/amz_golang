package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	portStr = ":9021"
)

func TcpNetServerLoop() {
	cm, err := net.Listen("tcp", portStr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("tcp server start at:%v\n", portStr)
	defer cm.Close()

	for {
		conn, err := cm.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go NetHandler(conn)
	}

}
func NetHandler(cn net.Conn) {
	for {
		data, err := bufio.NewReader(cn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		cn.Write([]byte("Hello " + data))
	}
}
func main() {
	//使用 socketio 模块
	// ServerIo()
	//使用原生net模块
	// TcpServerIo()

	//使用handler读取
	TcpNetServerLoop()
}
