package models

import (
	"github.com/gorilla/websocket"
)

// WsManager WebSocket 管理器
var WsManager = ClientManager{
	ClientGroup: make(map[string]map[string]*WsClient),
	Register:    make(chan *WsClient),
	UnRegister:  make(chan *WsClient),
	Broadcast:   make(chan *BoradcastData, 10),
}

// ClientManager websocket client Manager struct
type ClientManager struct {
	ClientGroup map[string]map[string]*WsClient
	Register    chan *WsClient
	UnRegister  chan *WsClient
	Broadcast   chan *BoradcastData
}

// boradcastData 广播数据
type BoradcastData struct {
	GroupID string
	Data    []byte
}

// wsClient Websocket 客户端
type WsClient struct {
	ID     string
	Group  string
	Socket *websocket.Conn
	Send   chan []byte
}
