package main

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	redigolib "github.com/gomodule/redigo/redis"
	"github.com/stvp/tempredis"
)

/*
红锁算法 redsync
1 它以毫秒为单位获取当前时间。
2 它尝试顺序获取所有 N 个实例中的锁，在所有实例中使用相同的键名和随机值。
在步骤 2 中，当在每个实例中设置锁时，客户端使用一个与锁自动释放总时间相比较小的超时来获取它。
例如，如果自动释放时间为 10 秒，则超时可能在 ~ 5-50 毫秒范围内。这可以防止客户端在尝试与已关闭的 Redis 节点通信时长时间保持阻塞：
3 如果一个实例不可用，我们应该尽快尝试与下一个实例通信。
客户端通过从当前时间中减去步骤 1 中获得的时间戳来计算获取锁所用的时间。
当且仅当客户端能够在大多数实例（至少 3 个）中获取锁时，且获取锁的总时间小于锁的有效时间，则认为锁已被获取。
4 如果获得了锁，则其有效时间被认为是初始有效时间减去经过的时间，如步骤 3 中计算的那样。
5 如果客户端由于某种原因未能获得锁（它无法锁定 N/2+1 个实例或有效时间为负数），它将尝试解锁所有实例（即使是它认为没有的实例）能够锁定）。
*/

/////////  LINUX 环境执行。 不支持windows
func main() {

	rsconfig := make(map[string]string)
	rsconfig["redis-server"] = "redis://192.168.30.129:7010/11"
	server, err := tempredis.Start(rsconfig) //tempredis.Config{})
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

	if err = mutex.Lock(); err != nil {
		panic(err)
	}

	if _, err = mutex.Unlock(); err != nil {
		panic(err)
	}
}
