package service

import (
	"mysockets/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// WsManager WebSocket 管理器
var WsManager = clientManager{models.WsManager}

// ClientManager websocket client Manager struct
type clientManager struct {
	models.ClientManager
}

// boradcastData 广播数据
type boradcastData struct {
	models.BoradcastData
}

// wsClient Websocket 客户端
type wsClient struct {
	models.WsClient
}

func (c *wsClient) Read() {
	// defer func() {
	// 	WsManager.UnRegister <- &c.WsClient
	// 	c.Socket.Close()
	// }()

	for {
		_, _, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *wsClient) Write() {
	// defer func() {
	// 	c.Socket.Close()
	// }()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.BinaryMessage, message)
		}
	}
}

// Start 启动 websocket 管理器
func (manager *clientManager) Start() {
	Logger.Printf("Websocket manage start")
	for {
		select {
		case client := <-manager.Register:
			Logger.Printf("Websocket client %s connect", client.ID)
			if manager.ClientGroup[client.Group] == nil {
				manager.ClientGroup[client.Group] = make(map[string]*models.WsClient)
			}
			manager.ClientGroup[client.Group][client.ID] = client
			Logger.Printf("Register client %s to %s group success", client.ID, client.Group)

		case client := <-manager.UnRegister:
			Logger.Printf("Unregister websocket client %s", client.ID)
			if _, ok := manager.ClientGroup[client.Group]; ok {
				if _, ok := manager.ClientGroup[client.Group][client.ID]; ok {
					close(client.Send)
					delete(manager.ClientGroup[client.Group], client.ID)
					Logger.Printf("Unregister websocket client %s from group %s success", client.ID, client.Group)

					if len(manager.ClientGroup[client.Group]) == 0 {
						Logger.Printf("Clear no client group %s", client.Group)
						delete(manager.ClientGroup, client.Group)
					}
				}
			}

		case data := <-manager.Broadcast:
			if groupMap, ok := manager.ClientGroup[data.GroupID]; ok {
				for _, conn := range groupMap {
					conn.Send <- data.Data
				}
			}
		}
	}
}

// RegisterClient 向 manage 中注册 client
func (manager *clientManager) RegisterClient(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Logger.Println("Error: websocket client connect %v error", ctx.Param("channel"))
		return
	}

	client := &wsClient{models.WsClient{
		ID:     uuid.NewV4().String(),
		Group:  ctx.Param("channel"),
		Socket: conn,
		Send:   make(chan []byte, 1024),
	}}

	manager.Register <- &client.WsClient
	go client.Read()
	go client.Write()
}

// Groupbroadcast 向指定的 Group 广播
func (manager *clientManager) Groupbroadcast(group string, message []byte) {
	data := &boradcastData{models.BoradcastData{
		GroupID: group,
		Data:    message,
	}}
	manager.Broadcast <- &data.BoradcastData
}
