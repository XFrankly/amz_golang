package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/trace"
	"golang.org/x/net/websocket"
)

var (
	port = flag.Int("port", 9080, "The server port")
)

type Event struct {
	// The fields of this struct must be exported so that the json module will be
	// able to write into them. Therefore we need field tags to specify the names
	// by which these fields go in the JSON representation of events.
	///// 这个结构的字段必须被导出，这样 json 模块才会被导出
	// 能够写入它们。因此我们需要字段标签来指定名称
	// 这些字段通过这些字段进入事件的 JSON 表示。
	X int `json:"x"`
	Y int `json:"y"`
}

// handleWebsocketEchoMessage handles the message e arriving on connection ws
// from the client.
func handleWebsocketEchoMessage(ws *websocket.Conn, e Event) error {
	// Log the request with net.Trace
	tr := trace.New("websocket.Receive", "receive")
	defer tr.Finish()
	tr.LazyPrintf("Got event %v\n", e)

	// Echo the event back as JSON ,只收Event x y {"x":12, "y":220, "z":21.2}
	err := websocket.JSON.Send(ws, e)
	if err != nil {
		return fmt.Errorf("Can't send: %s", err.Error())
	}
	return nil
}

// websocketEchoConnection handles a single websocket echo connection - ws.
// websocketEchoConnection 处理单个 websocket 回显连接 - ws。
func websocketEchoConnection(ws *websocket.Conn) {
	log.Printf("Client connected from %s", ws.RemoteAddr())
	for {
		var event Event
		//接受格式 {"hello":"world"}， {"x":2,"y":110}
		err := websocket.JSON.Receive(ws, &event)
		fmt.Printf("event rec:%+v, err:%+v\n", event, err)
		if err != nil {
			log.Printf("Receive failed: %s; closing connection...", err.Error())
			if err = ws.Close(); err != nil {
				log.Println("Error closing connection:", err.Error())
			}
			break
		} else {
			if err := handleWebsocketEchoMessage(ws, event); err != nil {
				log.Println(err.Error())
				break
			}
		}
	}
}

// websocketTimeConnection handles a single websocket time connection - ws.
// 处理单个 websocket 处理单个 websocket 时间连接
func websocketTimeConnection(ws *websocket.Conn) {
	for range time.Tick(1 * time.Second) {
		// Once a second, send a message (as a string) with the current time.
		websocket.Message.Send(ws, time.Now().Format("Mon, 02 Jan 2006 15:04:05 PST"))
	}
}

func main() {
	flag.Parse()
	// Set up websocket servers and static file server. In addition, we're using
	// net/trace for debugging - it will be available at /debug/requests.
	/// // 设置 websocket 服务器和静态文件服务器。此外，我们正在使用
	// 用于调试的 net/trace - 它将在 /debug/requests 中可用。
	http.Handle("/wsecho", websocket.Handler(websocketEchoConnection))
	http.Handle("/wstime", websocket.Handler(websocketTimeConnection))
	http.Handle("/", http.FileServer(http.Dir("static/html")))

	log.Printf("Server listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
