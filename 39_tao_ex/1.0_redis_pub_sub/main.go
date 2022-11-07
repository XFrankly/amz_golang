package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	logger         = log.New(os.Stderr, "INFO -", 13)
	ctx            = context.Background()
	InnerAddr      = "192.168.30.131:6379"
	InnerDB        = 0
	InnerSize      = 100
	InnerPasswd    = ""
	InnerPubChannl = "boards:apis:counts" // 广播 监控直播数据
	RedisInnerCfg  = map[string]string{
		"network":  "tcp",
		"address":  InnerAddr,
		"db":       strconv.Itoa(InnerDB),
		"passwrod": InnerPasswd,
		"poolSize": strconv.Itoa(InnerSize),
		"poolFIFO": "true",
	}
)

// 初始化redis 从内部配置
func InnerRedis(onlyMain bool) *redis.Client {
	var client *redis.Client
	if onlyMain {
		db_no, _ := strconv.Atoi(RedisInnerCfg["db"])
		client = redis.NewClient(&redis.Options{
			Addr:     RedisInnerCfg["address"],
			Password: RedisInnerCfg["password"], // no password set
			DB:       db_no,                     // use default DB
			PoolFIFO: true,
			PoolSize: 1000,
		})
		logger.Printf("INFO Redis connected:%v\n", client)
	}

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Println("ERROR redis connect ping failed, err:", err)
	} else {
		logger.Println("INFO Redis connected:", pong)
	}
	return client
}

func Do_debug() {
	client := InnerRedis(true)

	//為客戶端訂閱給定的模式。
	// 可以省略模式以創建空訂閱。
	pubsub := client.PSubscribe(ctx, InnerPubChannl)
	defer pubsub.Close()
	logger.Println("pubsub.String():", pubsub.String())
	// Expect(pubsub.String()).To(Equal("PubSub(mychannel*)"))
	info, err := pubsub.Receive(ctx)
	if err != nil {
		msg := fmt.Sprintf("err with pubsub.Receive:%v\n", err)
		panic(msg)
	}
	logger.Printf("pubsub Receive:%v\n", info)

	infos, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		msg := fmt.Sprintf("err with pubsub.Receive:%v\n", err)
		panic(msg)
	}
	logger.Printf("pubsub ReceiveMessage:%v\n", infos)

	for {
		// 写入 pub
		msg := fmt.Sprintf("%v,%v", "/stream/live/", time.Now().String())
		time.Sleep(time.Second * 1)
		i, err := client.Publish(ctx, InnerPubChannl, msg).Result()
		logger.Printf("Publish  ReceiveMessage:%v, err:%v \n", i, err)

		if err != nil {
			panic("publish err")
		}
		// 读取 sub
		select {
		/*
						Tick 是 NewTicker 的便捷包裝器，提供對滴答的訪問
			// 僅通道。
			雖然 Tick 對於不需要關閉的客戶端很有用
			// 股票代碼，請注意，沒有辦法將其關閉底層
			// Ticker 不能被垃圾收集器回收；它“洩漏”。
			// 與 NewTicker 不同，如果 d <= 0，Tick 將返回 nil
		*/
		//每毫秒检测一次
		case <-time.Tick(time.Microsecond * 1):
			//持续一分钟，如果一分钟内没有消息，则超时一次
			infosTimeOut, err := pubsub.ReceiveTimeout(ctx, time.Minute*1)
			if err != nil {
				msg := fmt.Sprintf("err with pubsub.ReceiveTimeout:%v\n", err)
				panic(msg)
			}
			logger.Printf("pubsub infosTimeOut:%v, publish:%v \n", infosTimeOut, i)
		}
	}

}

func main() {
	Do_debug()
}
