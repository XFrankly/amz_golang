package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
	interaddr := ws.RemoteAddr()
	// switch interaddr.(type) {
	// case string:

	// }
	fmt.Println("ws:", ws.MaxPayloadBytes, ws.PayloadType, ws.Config().Protocol)
	addr := fmt.Sprintf("%v", interaddr.String())
	fmt.Println("access from :", addr)
	ws.Write([]byte(addr))

}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
