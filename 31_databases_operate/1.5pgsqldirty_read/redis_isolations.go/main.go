package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Logg = log.New(os.Stderr, "INFO -:", 18)
)

////// 生成redlock的 随机字符串 key
// Bytes generates n random bytes
func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Base64 generates a random base64 string with length of n
func Base64(n int) string {
	return MakeString(n, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/")
}

// Base64 generates a random base62 string with length of n
func Base62(s int) string {
	return MakeString(s, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// Hex generates a random hex string with length of n
// e.g: 67aab2d956bd7cc621af22cfb169cba8
func Hex(n int) string { return hex.EncodeToString(Bytes(n)) }

// list of default letters that can be used to make a random string when calling String
// function with no letters provided
var defLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// String generates a random string using only letters provided in the letters parameter
// if user ommit letters parameters, this function will use defLetters instead
func MakeString(n int, letters ...string) string {
	var letterRunes []rune
	if len(letters) == 0 {
		letterRunes = defLetters
	} else {
		letterRunes = []rune(letters[0])
	}

	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(letterRunes))
	// on each loop, generate one random rune and append to output
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(Bytes(4))%l])
	}
	return bb.String()
}

//////////////////////////////////////////////////////////////////////////////
//////// 红锁的实现

const (
	// MinLockExpire 是最小的 锁过期时间。 MinLockExpire is the minimum lock expire time.
	MinLockExpire = 1000

	// 默认的 锁过期时间 DefaultLockExpire is the default lock expire time.
	DefaultLockExpire = 3000

	// 最大刷新次数 是最大重试次数 如果刷新连续失败，则自动刷新。
	// MaxRefreshRetryTimes is the maximum retry times of
	// auto refresh if the refresh fails continuously.
	MaxRefreshRetryTimes = 100
)

const (
	lockNameKeyFmt   = "__lock:%s__"
	lockSignalKeyFmt = "__signal:%s__"
)

var (
	/// 各场景的默认错误类型
	// ErrLockWithoutName 在创建具有空名称的锁时返回。 is returned when creating a lock with a empty name.
	ErrLockWithoutName = errors.New("empty lock name")

	// ErrLockExpireTooSmall 返回在 当创建一个过期时间小于 300 毫秒的锁时。 is returned when creating a lock with expire smaller than 300ms.
	ErrLockExpireTooSmall = errors.New("lock expiration too small")

	// ErrAlreadyAcquired 返回在 当尝试锁定已经获得的锁时。 is returned when try to lock an already acquired lock.
	ErrAlreadyAcquired = errors.New("lock already acquired")

	// ErrNotAcquired 返回在 当无法获取锁时。 is returned when a lock cannot be acquired.
	ErrNotAcquired = errors.New("lock not acquired")

	// ErrLockNotHeld 返回在 尝试释放未获取的锁时。 is returned when trying to release an unacquired lock.
	ErrLockNotHeld = errors.New("lock not held")
)

var (
	luaRefreshLock = redis.NewScript(`
		if redis.call("get", KEYS[1]) ~= ARGV[2] then
        	return 1
    	else
        	redis.call("pexpire", KEYS[1], ARGV[1])
        	return 0
    	end
	`)
	luaUnlock = redis.NewScript(`
		if redis.call("get", KEYS[1]) ~= ARGV[1] then
        	return 1
    	else
        	redis.call("del", KEYS[2])
			redis.call("lpush", KEYS[2], 1)
			redis.call("expire", KEYS[2], 600)
        	redis.call("del", KEYS[1])
        	return 0
    	end
	`)
)

//红锁代码 一个redis锁 结构体
// RedLock represents a redis lock.
type RedLock struct {
	cli                *redis.Client // redis client
	name               string        // redis key of lock
	holder             string        // lock holder name
	signalName         string        // redis key of lock release signal
	expiration         int           // expiration of lock in milliseconds
	autoRefresh        bool          // automatically refresh the lock or not
	refreshInterval    int           // refresh interval if autoRefresh is enabled
	failedRefreshCount int           // count of continuous failed auto refresh
	stopRefresh        chan struct{} // channel used to notify the background auto-refresh goroutine to stop
}

/*
/ New 创建并返回一个新的 redis 锁。
//
// 注意没有过期时间的分布式锁是危险的，所以总是需要过期时间。如果没有给出过期时间，或过期 <= 0，
// 将使用默认的 3 秒过期时间。
//
// 此外，一个非常小的过期时间，例如 5ms，根本没有意义。锁可能在调用者刚获取后就已经过期了
// 考虑调用者和redis服务器之间的网络通信时间。总是给出大于 1s 的过期时间。
// 启用了 autoRefresh 的锁将定期刷新其过期时间 后台 goroutine 让锁持有者可以一直持有锁
// 在释放它之前。刷新间隔 始终是锁过期的 2/3。
*/
func New(cli *redis.Client, name string, expiration int, autoRefresh bool) (*RedLock, error) {
	///redis v8客户端
	if name == "" {
		return nil, ErrLockWithoutName
	}
	if expiration <= 0 {
		expiration = DefaultLockExpire
	}

	if expiration < MinLockExpire {
		return nil, ErrLockExpireTooSmall
	}

	lock := &RedLock{
		cli:             cli,
		name:            fmt.Sprintf(lockNameKeyFmt, name),
		signalName:      fmt.Sprintf(lockSignalKeyFmt, name),
		holder:          MakeString(18),
		expiration:      expiration,
		autoRefresh:     autoRefresh,
		refreshInterval: int(float32(expiration) * 2 / 3),
	}
	return lock, nil
}

//// // hasAcquired 如果调用者已获取锁，则返回。
func (rd *RedLock) hasAcquired(ctx context.Context) (bool, error) {
	holder, err := rd.cli.Get(ctx, rd.name).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return holder == rd.holder, nil
}

func (rd *RedLock) refresh() error {
	timeout := time.Duration(rd.refreshInterval) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ret, err := luaRefreshLock.Run(ctx, rd.cli, []string{rd.name}, rd.expiration,
		rd.holder).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrLockNotHeld
		}
		return err
	}
	if i, ok := ret.(int64); !ok || i != 0 {
		return ErrLockNotHeld
	}
	return nil
}

//// // autoRefresh 在后台运行并自动刷新锁的过期时间。
func (rd *RedLock) startAutoRefresh() {
	rd.stopRefresh = make(chan struct{})
	go func() {
		interval := time.Duration(rd.refreshInterval) * time.Millisecond
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-rd.stopRefresh:
				return
			case <-ticker.C:
				err := rd.refresh()
				if err != nil {
					if errors.Is(err, ErrLockNotHeld) {
						return
					}
					rd.failedRefreshCount++
					if rd.failedRefreshCount > MaxRefreshRetryTimes {
						return
					}
					continue
				}
				rd.failedRefreshCount = 0
			}
		}
	}()
}

/// 锁住红锁 如果没有获取锁或触发上下文超时，则返回错误。
func (rd *RedLock) Lock(ctx context.Context, block bool) error {
	yes, err := rd.hasAcquired(ctx) // 是否已经获取了锁
	if err != nil {
		return err
	}
	if yes {
		return ErrAlreadyAcquired
	}
	yes, err = rd.cli.SetNX(ctx, rd.name, rd.holder, time.Duration(rd.expiration)*time.Millisecond).Result()

	if err != nil {
		return err
	}
	if yes {
		if rd.autoRefresh {
			rd.startAutoRefresh()
		}
		return nil
	}
	if !block {
		return ErrNotAcquired
	}

	var leftTime time.Duration
	deadline, ok := ctx.Deadline()
	if ok {
		leftTime = deadline.Sub(time.Now())
	}
	if leftTime <= 0 {
		return ErrNotAcquired
	}

	for {
		_, err := rd.cli.BLPop(ctx, leftTime, rd.signalName).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return ErrNotAcquired
			}
			return err
		}
		yes, err = rd.cli.SetNX(ctx, rd.name, rd.holder, time.Duration(rd.expiration)*time.Microsecond).Result()
		if err != nil {
			return err
		}
		if yes {
			return nil
		}
	}
}

///// 解锁 非锁定的红锁
func (rd *RedLock) Unlock(ctx context.Context) error {
	ret, err := luaUnlock.Run(ctx, rd.cli, []string{rd.name, rd.signalName}, rd.holder).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrLockNotHeld
		}
		return err
	}
	if i, ok := ret.(int64); !ok || i != 0 {
		return ErrLockNotHeld
	}
	if rd.autoRefresh {
		close(rd.stopRefresh)
	}
	return nil
}

////////////////////////////////////////////////////////////////////// 使用redlock
// const (
// 	Addr     = "192.168.30.131:6379"
// 	Password = ""  // no password set
// 	DBnumber = 0   // use default DB
// 	PoolSize = 100 // 连接池大小
// )

// type RedisCfg struct {
// 	Addr     string
// 	Password string
// 	DBNumber int
// 	PoolSize int
// }

// var (
// 	rdb  *redis.Client
// 	Rcfg = RedisCfg{Addr: Addr, Password: Password, DBNumber: DBnumber, PoolSize: PoolSize}
// 	Logg = log.New(os.Stderr, "Redis INFO -:", 18)
// 	ctx  = context.Background()
// )

const (
	Addr     = "192.168.30.129:7010"
	Password = ""  // no password set
	DBnumber = 0   // use default DB
	PoolSize = 100 // 连接池大小
)

type RedisCfg struct {
	Addr     string
	Password string
	DBNumber int
	PoolSize int
}

var (
	Rdb  *redis.Client
	Rcfg = RedisCfg{Addr: Addr, Password: Password, DBNumber: DBnumber, PoolSize: PoolSize}

	Ctx = context.Background()
)

// 初始化连接
func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Rcfg.Addr,
		Password: Rcfg.Password, // no password set
		DB:       Rcfg.DBNumber, // use default DB
		PoolSize: Rcfg.PoolSize, // 连接池大小
	})

	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	Logg.Println("new context:", ctx1, cancel)
	defer cancel()

	testrdb, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		msg := fmt.Sprintf("init rdb %v fail %+v\n.", testrdb, err)
		panic(msg)
	}

}
func main() {

	Logg.Println(Base64(12), MakeString(50))
	//// redlock 名称 和 超时时间 5000 毫秒
	newRedLock, err := New(Rdb, MakeString(15), 5000, true)
	Logg.Println("get new red lock, err", newRedLock, err)
	if err != nil {
		Logg.Println("fail to new red lock.")
		panic(err)
	}
	err2 := newRedLock.Lock(Ctx, true)
	Logg.Println("lock err", err2)
	if err2 != nil {

		panic(err)
	}

	unerr := newRedLock.Unlock(Ctx)
	Logg.Println("unlock err", unerr)
	if unerr != nil {

		panic(unerr)
	}

}
