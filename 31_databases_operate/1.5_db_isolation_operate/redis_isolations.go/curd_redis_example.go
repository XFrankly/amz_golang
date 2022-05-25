// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	redis "github.com/go-redis/redis/v8"
// )

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
// 	Rdb  *redis.Client
// 	Rcfg = RedisCfg{Addr: Addr, Password: Password, DBNumber: DBnumber, PoolSize: PoolSize}
// 	Logg = log.New(os.Stderr, "Redis INFO -:", 18)
// 	Ctx  = context.Background()
// )

// // 初始化连接
// func init() {
// 	Rdb = redis.NewClient(&redis.Options{
// 		Addr:     Rcfg.Addr,
// 		Password: Rcfg.Password, // no password set
// 		DB:       Rcfg.DBNumber, // use default DB
// 		PoolSize: Rcfg.PoolSize, // 连接池大小
// 	})

// 	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	Logg.Println("new context:", ctx1, cancel)
// 	defer cancel()

// 	testrdb, err := Rdb.Ping(Ctx).Result()
// 	if err != nil {
// 		msg := fmt.Sprintf("init rdb %v fail %+v\n.", testrdb, err)
// 		panic(msg)
// 	}

// }

// /*context
// 			// ctx backgroumd
// 			背景返回一个非零的空上下文。它永远不会被取消，没有值，
// 			并且没有截止日期。它通常由 main 函数使用， 初始化和测试，
// 		并作为传入的顶级上下文 要求。
// 		上下文携带截止日期、取消信号和其他值 API 边界。
// // Context 的方法可以被多个 goroutine 同时调用。
// 	# ctx 属性1 deadline
// 		//Deadline
// 	/ Deadline 返回代表此上下文完成工作的时间
// // 应该取消。没有截止日期时，截止日期返回 ok==false
// // 放。对 Deadline 的连续调用返回相同的结果。
// 	# ctx属性2 Done()
// 	// Done 返回一个在代表 当前对象 完成工作时，在关闭的通道channel 上下文应该被取消。
// 	如果上下文可以，Done 可能会返回 nil 并且永远不会被取消。
// 	对 Done 的连续调用返回相同的值。
// // Done 通道的关闭可能是异步发生的，比如取消函数返回后。
// //
// // WithCancel 安排在调用 cancel 时关闭 Done；
// // WithDeadline 安排 Done 在截止日期时关闭过期；
// 	WithTimeout 安排在超时时关闭 Done 。
// Done 用于在 select 语句中使用：
// //
// // // Stream函数 使用 DoSomething 生成值并将它们发送到 out 直到 DoSomething 返回错误或 ctx.Done 关闭。
// func Stream(ctx context.Context, out chan<- Value) error {
// 	//  	for {
// 	//  		v, err := DoSomething(ctx)
// 	//  		if err != nil {
// 	//  			return err
// 	//  		}
// 	//  		select {
// 	//  		case <-ctx.Done():
// 	//  			return ctx.Err()
// 	//  		case out <- v:
// 	//  		}
// 	//  	}
// 	//  }
// 	有关如何使用的更多示例，请参见 https://blog.golang.org/pipelines  Done用于取消  通道。

// 	属性3 Err()
// 	// 如果 Done 尚未关闭，则 Err 返回 nil。
// // 如果 Done 关闭，Err 返回一个非 nil 错误，解释原因：
// 		如果上下文被取消则取消
// 		如果上下文的最后期限已过，则为 DeadlineExceeded。
// 		Err 返回非 nil 错误后，对 Err 的连续调用返回相同的错误。

// 	属性4 Value()
// 	 Value 返回与此上下文关联的值作为 key，或 nil 如果没有值与键关联。
// 	 连续调用 Value with 相同的键返回相同的结果。

// // 仅对传输的请求范围数据使用上下文值 进程和 API 边界，不用于将可选参数传递给职能。
// // 键标识上下文中的特定值。希望的功能在 Context 中存储值通常在全局中分配一个键变量然后使用该键作为 context.WithValue 的参数和
// // 上下文。
// 值。键可以是任何支持相等的类型； 包应该将键定义为未导出的类型以避免 碰撞。
// // 定义 Context 键的包应该提供类型安全的访问器 对于使用该键存储的值：
// // 包用户定义了存储在上下文中的用户类型。 Package user defines a User type that's stored in Contexts.
// 	// 	package user
// 	//
// 	// 	import "context"
// 	//
// 	// 	// User 是存储在上下文中的值的类型。 User is the type of value stored in the Contexts.
// 	// 	type User struct {...}
// 	//key 是此包中定义的键的未导出类型。 这可以防止与其他包中定义的键发生冲突。
// 	// 	// key is an unexported type for keys defined in this package.
// 	// 	// This prevents collisions with keys defined in other packages.
// 	// 	type key int
// 	//userKey 是 Contexts 中 user.User 值的键。这是未导出；
// 	客户端使用 user.NewContext 和 user.FromContext  而不是直接使用这个键。
// 	// 	// userKey is the key for user.User values in Contexts. It is
// 	// 	// unexported; clients use user.NewContext and user.FromContext
// 	// 	// instead of using this key directly.
// 	// 	var userKey key
// 	//
// 	// 	// NewContext 返回一个带有值 u 的新 Context。 NewContext returns a new Context that carries value u.
// 	// 	func NewContext(ctx context.Context, u *User) context.Context {
// 	// 		return context.WithValue(ctx, userKey, u)
// 	// 	}
// 	//
// 	// 	// FromContext 返回存储在 ctx 中的用户值（如果有任何）。 FromContext returns the User value stored in ctx, if any.
// 	// 	func FromContext(ctx context.Context) (*User, bool) {
// 	// 		u, ok := ctx.Value(userKey).(*User)
// 	// 		return u, ok
// 	// 	}

// */

// func V8Example() {

// 	err := Rdb.Set(Ctx, "keyredis", "valueredis", 0).Err()
// 	// Logg.Println("ctx info of key:", ctx)
// 	// f := ctx.Value //("keyredis") = "valueredis"
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := Rdb.Get(Ctx, "keyredis").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("keyredis:", val)

// 	val2, err := Rdb.Get(Ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// 	Logg.Printf("ctx:%+v\n", Ctx)
// }

// // set/get示例
// func redisExample() {
// 	err := Rdb.Set(Ctx, "score", 100, 0).Err()
// 	if err != nil {
// 		fmt.Printf("set score failed, err:%v\n", err)
// 		return
// 	}

// 	val, err := Rdb.Get(Ctx, "score").Result()
// 	if err != nil {
// 		fmt.Printf("get score failed, err:%v\n", err)
// 		return
// 	}
// 	fmt.Println("score", val)

// 	val2, err := Rdb.Get(Ctx, "name").Result()
// 	if err == redis.Nil {
// 		fmt.Println("name does not exist")
// 	} else if err != nil {
// 		fmt.Printf("get name failed, err:%v\n", err)
// 		return
// 	} else {
// 		fmt.Println("name", val2)
// 	}
// }

// // zset示例
// func redisExample2() {
// 	zsetKey := "language_rank"
// 	languages := []*redis.Z{
// 		&redis.Z{Score: 90.0, Member: "Golang"},
// 		&redis.Z{Score: 98.0, Member: "Java"},
// 		&redis.Z{Score: 95.0, Member: "Python"},
// 		&redis.Z{Score: 97.0, Member: "JavaScript"},
// 		&redis.Z{Score: 99.0, Member: "C/C++"},
// 	}
// 	// ZADD
// 	num, err := Rdb.ZAdd(Ctx, zsetKey, languages...).Result()
// 	if err != nil {
// 		fmt.Printf("zadd failed, err:%v\n", err)
// 		return
// 	}
// 	fmt.Printf("zadd %d succ.\n", num)

// 	// 把Golang的分数加10
// 	newScore, err := Rdb.ZIncrBy(Ctx, zsetKey, 10.0, "Golang").Result()
// 	if err != nil {
// 		fmt.Printf("zincrby failed, err:%v\n", err)
// 		return
// 	}
// 	fmt.Printf("Golang's score is %f now.\n", newScore)

// 	// 取分数最高的3个
// 	ret, err := Rdb.ZRevRangeWithScores(Ctx, zsetKey, 0, 2).Result()
// 	if err != nil {
// 		fmt.Printf("zrevrange failed, err:%v\n", err)
// 		return
// 	}
// 	for _, z := range ret {
// 		fmt.Println(z.Member, z.Score)
// 	}

// 	// 取95~100分的
// 	op := &redis.ZRangeBy{
// 		Min: "95",
// 		Max: "100",
// 	}
// 	ret, err = Rdb.ZRangeByScoreWithScores(Ctx, zsetKey, op).Result()
// 	if err != nil {
// 		fmt.Printf("zrangebyscore failed, err:%v\n", err)
// 		return
// 	}
// 	for _, z := range ret {
// 		fmt.Println(z.Member, z.Score)
// 	}
// }

// // 按通配符删除key
// func exampleTongPeiFu() {
// 	ctx := context.Background()
// 	iter := Rdb.Scan(ctx, 0, "prefix*", 0).Iterator()
// 	for iter.Next(ctx) {
// 		err := Rdb.Del(ctx, iter.Val()).Err()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	if err := iter.Err(); err != nil {
// 		panic(err)
// 	}
// }

// // func main() {
// // 	V8Example()

// // 	Logg.Printf("final ctx: %+v\n", <-ctx.Done()) //.Value("keyredis"))
// // }
