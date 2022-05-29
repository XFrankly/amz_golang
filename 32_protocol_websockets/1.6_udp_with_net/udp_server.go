package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
該函數將使用導入的包創建一個 UDP 服務器。

main() 函數在 arguments 變量中收集命令行參數並包括錯誤處理。
net.ListenUDP() 函數告訴應用程序偵聽傳入的 UDP 連接，這些連接在 for 循環中提供服務。
這是使程序成為 UDP 服務器的函數調用。
ReadFromUDP() 和 WriteToUDP() 函數分別用於從 UDP 連接讀取數據和將數據寫入 UDP 連接。
字節切片存儲在數據變量中，用於寫入所需的數據。
緩衝區變量還存儲一個字節切片，用於讀取數據。
由於 UDP 是無狀態協議，因此每個 UDP 客戶端都會得到服務，然後連接會自動關閉。
UDP 服務器程序只有在收到來自 UDP 客戶端的 STOP 關鍵字時才會退出。否則，服務器程序將繼續等待來自其他客戶端的更多 UDP 連接。

*/
var Logg = log.New(os.Stderr, "INFO -:", 18)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	Logg.Println("rand int:", random(1, 100))
	arguments := os.Args
	if len(arguments) == 1 {
		Logg.Println("Please give a port number")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		Logg.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		Logg.Println(err)
		return
	}
	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		Logg.Print("->", string(buffer[0:n-1]))
		receive_data := strings.TrimSpace(string(buffer[0:n]))
		if receive_data == "STOP" {
			Logg.Println("EXiting UDP server!")
			return
		}
		data := []byte(strconv.Itoa(random(1, 1001)))

		response := fmt.Sprintf("%v.%v", receive_data, string(data))
		Logg.Println("response with rand data:", string(response))
		dd, err2 := connection.WriteToUDP([]byte(response), addr)
		Logg.Println("write to udp", dd, err2)
		if err2 != nil {
			Logg.Println(err)
			return
		}
	}
}
