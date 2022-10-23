package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

func ProfucerOne() {
	conn, err := amqp.Dial(Config.AMQPConnRabbitMQ)
	HandleError(err, "Cannot connect")
	defer conn.Close()

	ampqChannel, err := conn.Channel()
	HandleError(err, "Cannot create amqp channel")
	defer ampqChannel.Close()

	queue, err := ampqChannel.QueueDeclare("add", true, false, false, false, nil)
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

func Consumer() {
	conn, err := amqp.Dial(Config.AMQPConnRabbitMQ)
	HandleError(err, "Cannot connect")
	defer conn.Close()

	ampqChannel, err := conn.Channel()
	HandleError(err, "Cannot create amqp channel")
	defer ampqChannel.Close()

	queue, err := ampqChannel.QueueDeclare("add", true, false, false, false, nil)
	HandleError(err, "couldn't declare add queue")

	err = ampqChannel.Qos(1, 0, false)
	HandleError(err, "could not config Qos")

	messageChannel, err := ampqChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	HandleError(err, "cannot register consume")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Our consumer read")
		for d := range messageChannel {
			log.Printf("Received a message:%s", d.Body)
			addTask := &AddTask{}

			err := json.Unmarshal(d.Body, addTask)
			if err != nil {
				log.Printf("Error in Decoding Json %s error", err)
			}
			log.Printf("Result of %d + %d is : %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			if err := d.Ack(false); err != nil {
				log.Printf("Error message:%s", err)
			} else {
				log.Printf("Ack message")
			}

		}
	}()
	<-stopChan
}

func main() {
	// 生产一个队列
	ProfucerOne()

	//消费队列
	Consumer()

}
