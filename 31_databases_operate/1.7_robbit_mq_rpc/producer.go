package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

type Configuration struct {
	AMQPConnectionURL string
	AMQPConnRabbitMQ  string
}

type AddTask struct {
	Number1 int
	Number2 int
}

var Config = Configuration{
	// AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
	AMQPConnRabbitMQ: "amqp://user:bitnami@192.168.30.130:5672/",
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(Config.AMQPConnRabbitMQ)
	HandleError(err, "Cannot connect")
	defer conn.Close()

	ampqChannel, err := conn.Channel()
	HandleError(err, "Cannot create amqp channel")
	defer ampqChannel.Close()

	queue, err := ampqChannel.QueueDeclare("task_queue", true, false, false, false, nil)
	HandleError(err, "couldn't declare add queue")

	rand.Seed(time.Now().UnixNano())

	addTask := AddTask{
		Number1: rand.Intn(999),
		Number2: rand.Intn(999),
	}
	body, err := json.Marshal(addTask)
	if err != nil {
		HandleError(err, "Error encoding JSON")
	}

	err = ampqChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("error in publishing msg %s", err)
	}

	log.Printf("AddTask:%d+%d", addTask.Number1, addTask.Number2)
}
