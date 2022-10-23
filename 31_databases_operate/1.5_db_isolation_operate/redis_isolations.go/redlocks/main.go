package main

import (
	"log"
	"os"
	redigolib "redlock/redis"
	"redlock/redsync"
	"redlock/redsync/redis/redigo"
	"redlock/tempredis"
	"time"
)

var (
	Logg = log.New(os.Stderr, "INFO -:", 18)
)

/////////  LINUX 环境执行。 不支持windows
func main() {
	rsconfig := make(map[string]string)
	rsconfig["redis-server"] = "redis://192.168.30.129:7010/0"
	tempConf := tempredis.Config{}
	Logg.Printf("get config from os:%+v\n", tempConf)
	server, err := tempredis.Start(tempConf) //tempredis.Config{})
	Logg.Printf("tempredis server:%+v\n", server)
	Logg.Printf("%+v\n", err)
	if err != nil {
		panic(err)
	}
	defer server.Term()

	pool := redigo.NewPool(&redigolib.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigolib.Conn, error) {
			return redigolib.Dial("unix", server.Socket())
		},
		TestOnBorrow: func(c redigolib.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	})

	rs := redsync.New(pool)

	mutex := rs.NewMutex("test-redsync")
	Logg.Printf("redsync pool rs:%+v\n", rs)
	Logg.Printf("pool:%+v\n", pool)
	Logg.Printf("mutex:%+v\n", mutex)
	if err = mutex.Lock(); err != nil {
		panic(err)
	}

	if _, err = mutex.Unlock(); err != nil {
		panic(err)
	}
}
