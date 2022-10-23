package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx = context.Background()
)

type redisBuf struct {
	*redis.Client
	buf  *bufio.Reader
	conn net.Conn
}

func CacheRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.30.131:6379",
		Password: "",
		DB:       "0:", //cfg.DB,
		// MaxRetries:   cfg.MaxRetries,
		// PoolSize:     cfg.PoolSize,
		// MinIdleConns: cfg.MinIdleConns,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis err:%v", err)
	}

	return client, nil
}

func setup(t *testing.T, option redis.Options) (*pipe, *redisBuf, func(), func()) {
	if option.CacheSizeEachConn <= 0 {
		option.CacheSizeEachConn = DefaultCacheBytes
	}
	n1, n2 := net.Pipe()
	mock := &redisBuf{
		buf:  bufio.NewReader(n2),
		conn: n2,
	}
	go func() {
		mock.Expect("HELLO", "3").
			Reply(RedisMessage{
				typ: '%',
				values: []RedisMessage{
					{typ: '+', string: "version"},
					{typ: '+', string: "6.0.0"},
				},
			})
		if !option.DisableCache {
			mock.Expect("CLIENT", "TRACKING", "ON", "OPTIN").
				ReplyString("OK")
		}
	}()
	p, err := newPipe(func() (net.Conn, error) { return n1, nil }, &option)
	if err != nil {
		t.Fatalf("pipe setup failed: %v", err)
	}
	if info := p.Info(); info["version"].string != "6.0.0" {
		t.Fatalf("pipe setup failed, unexpected hello response: %v", p.Info())
	}
	return p, mock, func() {
			go func() { mock.Expect("QUIT").ReplyString("OK") }()
			p.Close()
			mock.Close()
		}, func() {
			n1.Close()
			n2.Close()
		}
}

func TestWriteMultiPipelineFlush(t *testing.T) {
	p, mock, cancel, _ := setup(t, redis.ClientOption{})
	defer cancel()
	times := 2000
	wg := sync.WaitGroup{}
	wg.Add(times)

	for i := 0; i < times; i++ {
		go func() {
			for _, resp := range p.DoMulti(context.Background(), cmds.NewCompleted([]string{"PING"}), cmds.NewCompleted([]string{"PING"})) {
				ExpectOK(t, resp)
			}
		}()
	}

	for i := 0; i < times; i++ {
		mock.Expect("PING").Expect("PING").ReplyString("OK").ReplyString("OK")
	}
}
