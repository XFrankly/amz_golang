package main  
import (
	"log"
	"runtime"
	"socketexample/gosocketio"
	"socketexample/gosocketio/models"
	"socketexample/gosocketio/models/protocol"
	"socketexample/gosocketio/models/transport"
	"time"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func websocketClients() {
	/*
			设置可以执行的最大 CPU 数
		// 同时返回之前的设置。它默认为
		// runtime.NumCPU 的值。如果 n < 1，则不会更改当前设置。
		// 当调度器改进时，这个调用将消失。
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 3811, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *models.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(models.OnDisconnection, func(h *models.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(models.OnConnection, func(h *models.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)

	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")
}
func main() {
	log.Println(protocol.MessageTypeOpen)
	log.Println(protocol.MessageTypeClose)
	log.Println(protocol.MessageTypePing)
	log.Println(protocol.MessageTypePong)
	log.Println(protocol.MessageTypeEmpty)
	log.Println(protocol.MessageTypeEmit)

	log.Println(protocol.MessageTypeAckRequest)
	log.Println(protocol.MessageTypeAckResponse)
	websocketClients()
}
