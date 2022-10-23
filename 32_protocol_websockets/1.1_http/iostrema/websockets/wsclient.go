package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := ws.Write([]byte(fmt.Sprintf("hello, world! AccessNo.%d\n", i))); err != nil {
			log.Fatal(err)
		}
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received: %s.\n", msg[:n])
		time.Sleep(time.Second * 1)
	}

}
