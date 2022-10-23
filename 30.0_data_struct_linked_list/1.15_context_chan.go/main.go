package main

import (
	"context"
	"log"
	"os"
	"reflect"
	"time"
)

/*
context.Context
属性
	(c *cancelCtx) Value()
	(c *cancelCtx) Done()
	(c *cancelCtx) Err()

	(c *cancelCtx) String()

	/// 创建截止期限 context
	 WithDeadline(parent Context, d time.Time) (Context, CancelFunc)

	 (c *timerCtx) Deadline() (deadline time.Time, ok bool)
	//// 创建一个有时限的 context
	 WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

	 //// 创建一个有Value值的 context
	 WithValue(parent Context, key, val any) Context


上下文管理的几个例子
example 1
这个例子演示了使用可取消上下文来防止 goroutine 泄漏。
在示例函数结束时，由 gen 启动的 goroutine 将返回而不会泄漏。



*/
const (
	shortDuration = 10 * time.Millisecond
)

var (
	Logg = log.New(os.Stderr, "context INFO -:", 18)
)

func Example1() {

	/*
		// gen 在一个单独的 goroutine 中生成整数，并且将它们发送到返回的通道。
		// gen 的调用者需要取消一次上下文 ，保证他们完成消费生成的整数不泄漏
		// 由 gen 启动的内部 goroutine。
	*/

	gen := func(ctx context.Context) <-chan int {
		//// 生成器
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return /// 返回将不泄露 goroutine 。 leak the  goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	/// 获取一个 自动取消 的 上下文管理器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()                      // 当我们完成int的消费，取消
	Logg.Printf("func gen:%+v\n", &gen) /// 一个地址
	for n := range gen(ctx) {
		Logg.Println(n)
		if n == 5 {
			break
		}
	}
}

func Example2() {
	/*
			WithDeadline 返回父上下文的副本，截止日期调整为不迟于 d。
			如果父节点的截止日期已经早于 d，WithDeadline(parent, d) 在语义上等价于父节点。
			返回的上下文的 Done 通道在截止日期到期、调用返回的取消函数或父上下文的 Done 通道关闭时关闭，以先发生者为准。

		取消此上下文会释放与其关联的资源，因此代码应在此上下文中运行的操作完成后立即调用取消。

		这个例子传递了一个带有任意截止日期的上下文来告诉一个阻塞函数它应该在它到达它时立即放弃它的工作。
	*/

	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	//// // 即使 ctx 会过期，最好调用它的 任何情况下的取消函数。
	// 否则可能会保留 上下文及其父级的存活时间超过了必要的时间。
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		Logg.Println("ctx overslept.")
	case <-ctx.Done(): /// 让 ctx 关闭
		Logg.Println("ctx err:", ctx.Err()) ////prints "context deadline exceeded"
	}
}

//// WithTimeOut, 实际上是返回了一个 WithDeadline函数的对象
func Example3() {
	/*
			WithTimeout 返回 WithDeadline(parent, time.Now().Add(timeout))。

		取消此上下文会释放与其关联的资源，
		因此代码应在此上下文中运行的操作完成后立即调用取消：
	*/
	/// 这个例子传递了一个带有超时的上下文来告诉一个阻塞函数它应该在超时后放弃它的工作。
	// 传递一个带有超时的上下文来告诉阻塞函数它 超时后应该放弃它的工作。
	ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		Logg.Println("overslept")
	case <-ctx.Done(): //// 制作超时
		Logg.Println("ctx err :", ctx.Err()) /// prints "context deadline exceeded"
	}
}

func Example4() {
	/*
			WithValue 返回一个 parent 的副本，其中与 key 关联的值为 val。

		仅将上下文值用于传输流程和 API 的请求范围数据，而不用于将可选参数传递给函数。

		提供的键必须是可比较的，并且不应该是字符串类型或任何其他内置类型，以避免使用上下文的包之间发生冲突。
		 WithValue 的用户应该为键定义自己的类型。
		为避免在分配给 interface{} 时进行分配，上下文键通常具有具体类型 struct{}。
		或者，导出的上下文键变量的静态类型应该是指针或接口。
	*/
	/// 这个例子演示了如何将一个值传递给上下文，以及如何在它存在时检索它
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil { //// 查询ctx 的key 的对应值
			Logg.Println("found value:", v, "from k:", k, "k type:", reflect.TypeOf(k).Name())
			return
		}
		Logg.Println("key not found:", k)
	}
	k := favContextKey("language.")                         //// 创建一个 自定义 字符串 结构体，值为 language
	ctx := context.WithValue(context.Background(), k, "Go") /// 创建一个上下文，k为自定义字符串结构体，值为 Go

	f(ctx, k)                      //// 尝试 查找自定义字符串结构体 language，将打印找到
	f(ctx, favContextKey("color")) ///  尝试查找 自定义结构体 color， 打印 未能找到

	//// 创建一个
	ctx2 := context.WithValue(context.Background(), "key", "self value")

	f(ctx2, "key")

	//// 创建一个超时上下文
	ctx3, cancel3 := context.WithTimeout(context.Background(), shortDuration)
	Logg.Println("err1:", ctx3.Err()) ///  打印 nil, 此时还没有错误。

	time.Sleep(1 * time.Second)
	Logg.Println("ctx3 state:", ctx3, cancel3)
	Logg.Println(ctx3.Err()) ///  此时已经超时，打印 context deadline exceeded

	defer cancel3()
	Logg.Printf("%+v\n", ctx3.Done())
	Logg.Println(ctx3.Value("key")) /// 没有Value 返回nil
}

func NestedContext(parent context.Context) context.Context { /// 返回的接口类型
	return &Cont2{Cont1{"me"}, parent, "cont1", "ok"}
}

type Cont1 struct {
	Name string
}

type Cont2 struct {
	//// 接口类型
	Cont1
	context.Context
	key, val string
}

func main() {
	Example1()
	Example2()
	Example3()
	Example4()
}
