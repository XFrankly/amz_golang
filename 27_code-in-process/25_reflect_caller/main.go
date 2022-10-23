package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type Caller struct {
	Func        reflect.Value
	Args        reflect.Type
	ArgsPresent bool
	Out         bool
}

type Channel struct {
	Out chan string

	Alive             bool
	AliveLock         sync.Mutex
	ResultWaitersLock sync.RWMutex

	ClientIp          string
	HttpRequestHeader http.Header
}

var (
	ErrorCallerNotFunc     = errors.New("f is not function")
	ErrorCallerNot2Args    = errors.New("f should have 1 or 2 args")
	ErrorCallerMaxOneValue = errors.New("f should return not more than one value")
)

func NewChannel() *Channel {

	cc := make(chan string, 10)
	// return &Channel{
	// 	Out:      cc,
	// 	Alive:    true,
	// 	ClientIp: "127.0.0.1",
	// }
	return &Channel{
		// Conn:(*views.WebsocketConnection)(0xc00004e0b0),
		Out: cc,
		// Header:models.Header{Sid:"8LIAL5doOT3yUwwrYV-k", Upgrades:[]string{}, PingInterval:30000, PingTimeout:60000},
		Alive:     true,
		AliveLock: sync.Mutex{},
		// Ackproce:models.AckProcessor{Counter:0, CounterLock:sync.Mutex{state:0, sema:0x0},
		// ResultWaiters:map[int]chan string{},
		ResultWaitersLock: sync.RWMutex{},

		HttpRequestHeader: http.Header(nil)}

}

func (ca *Channel) SendMsg(s string) chan string {
	ca.AliveLock.Lock()
	defer ca.AliveLock.Unlock()
	if s == "" {
		s = "Hello"
	}
	ca.Out <- s
	return ca.Out
}

func (ca *Channel) GetMsg() string {
	ca.AliveLock.Lock()
	defer ca.AliveLock.Unlock()
	return <-ca.Out
}

func (ca *Channel) Close() {
	ca.Alive = false
}

/**
解析使用反射传递的函数，并存储其表示进一步调用消息或确认
*/
func NewCaller(f interface{}) (*Caller, error) {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		return nil, ErrorCallerNotFunc
	}

	fType := fVal.Type()
	if fType.NumOut() > 1 {
		return nil, ErrorCallerMaxOneValue
	}

	curCaller := &Caller{
		Func: fVal,
		Out:  fType.NumOut() == 1}

	if fType.NumIn() == 1 {
		curCaller.Args = nil
		curCaller.ArgsPresent = false
	} else if fType.NumIn() == 2 {
		curCaller.Args = fType.In(1)
		curCaller.ArgsPresent = true
	} else {
		return nil, ErrorCallerNot2Args
	}

	return curCaller, nil
}

/**
使用反射返回函数参数，因为它存在于其中
*/
func (c *Caller) GetArgs() interface{} {
	ca := c.Args
	fmt.Printf("get Args:%#v\n", ca)
	return reflect.New(ca).Interface()
}
func (ce *Caller) GetFunc() reflect.Value {
	return ce.Func
}

/**
使用反射从其表示中调用具有给定参数的函数
*/
func (c *Caller) CallFunc(h *Channel, args interface{}) []reflect.Value {
	//nil 是无类型的，所以使用正确类型的默认空值
	if args == nil {
		args = c.GetArgs()
	}
	//Elem 返回接口 v 包含的值
	// 或者指針 v 指向的那個。
	// 如果 v 的 Kind 不是接口或指針，它會恐慌。
	// 如果 v 為 nil，則返回零值。
	// a := []reflect.Value{reflect.ValueOf(h), reflect.ValueOf(args).Elem()}
	// if !c.ArgsPresent {
	// 	a = a[0:1]
	// }
	//Call 使用输入参数 in 调用函数 v。例如，如果 len(in) == 3，则 v.Call(in) 表示 Go 调用 v(in[0], in[1], in[2]) .
	//如果 v 的 Kind 不是 Func，则调用恐慌。它将输出结果作为值返回。
	//与 Go 一样，每个输入参数都必须可分配给函数对应输入参数的类型。
	//如果 v 是可变参数函数，则 Call 创建可变参数切片参数本身，并复制相应的值。
	fmt.Printf("c.Func:%#v, reflect Value:%#v\n", c.Func, reflect.TypeOf(c.Func).Name())
	fmt.Printf("func Kind:%#v \n", c.Func.Kind())
	//ValueOf函数返回一个Value类型值，该值代表运行时的数据

	var val reflect.Value

	//如果是golang 数据类型，直接返回值，否则需要使用 Elem()
	fmt.Printf("args all elem:%#v\n", reflect.ValueOf(args))

	if reflect.TypeOf(args).Name() == reflect.TypeOf(TestArgs{}).Name() {
		//Args结构体
		targs := args.(TestArgs)
		val = reflect.ValueOf(targs.Name)
	} else if reflect.TypeOf(args).Name() == reflect.TypeOf("a").Name() {
		//字符
		val = reflect.ValueOf(args.(string))
	} else if reflect.TypeOf(args).Kind() == reflect.TypeOf(struct{}{}).Kind() {
		//匿名结构体
		fmt.Println("its a struct{}{}")
	} else if reflect.TypeOf(args).Kind() == reflect.TypeOf(&struct{}{}).Kind() {
		//匿名指针
		fmt.Println("its a ptr &struct{}{}")

	} else {
		//未知场景
		panic(args)
		val = reflect.ValueOf("Hello, Whats That?")

	}
	fmt.Println("Call with args:", val)
	return c.Func.Call([]reflect.Value{val})
}

func StringInterface() interface{} {
	return struct{ s string }{"Worlds,"}
}

func CallFuncByReflect(s string) {
	chan1 := NewChannel()
	c1, err := NewCaller(chan1.SendMsg)
	c1.Args = reflect.TypeOf(chan1)
	fmt.Println("New caller:", c1, err)
	fmt.Println("caller:", c1.Args, c1.Out)

	//反射调用
	/*
				c:Caller{Func:reflect.Value{typ:(*reflect.rtype)(0x67d560),
					ptr:(unsafe.Pointer)(0x6ecb28), flag:0x13},
					Args:reflect.Type(nil),
					ArgsPresent:false,
					Out:false}
			   h:&models.Channel{Conn:(*views.WebsocketConnection)(0xc0000a6080),
				Out:(chan string)(0xc00005a180),
				Header:models.Header{Sid:"_1D0XpMMMpZUWaiMX8mm",
				Upgrades:[]string{},
				PingInterval:30000,
				PingTimeout:60000},
				Alive:true,
				AliveLock:sync.Mutex{state:0, sema:0x0},
				Ackproce:models.AckProcessor{Counter:0,
											CounterLock:sync.Mutex{state:0, sema:0x0},
											ResultWaiters:map[int]chan string{},
											ResultWaitersLock:sync.RWMutex{
												w:sync.Mutex{state:0, sema:0x0},
												writerSem:0x0,
												readerSem:0x0,
												readerCount:0,
												readerWait:0}},
											ClientIp:"",
											HttpRequestHeader:http.Header(nil)},
			args:&struct {}{}

		&models.Channel{
			Conn:(*views.WebsocketConnection)(0xc00004e0b0),
			Out:(chan string)(0xc00017e120),
			Header:models.Header{Sid:"8LIAL5doOT3yUwwrYV-k", Upgrades:[]string{}, PingInterval:30000, PingTimeout:60000},
			Alive:true, AliveLock:sync.Mutex{state:0, sema:0x0},
			Ackproce:models.AckProcessor{Counter:0, CounterLock:sync.Mutex{state:0, sema:0x0},
			ResultWaiters:map[int]chan string{},
			ResultWaitersLock:sync.RWMutex{w:sync.Mutex{state:0, sema:0x0}, writerSem:0x0, readerSem:0x0, readerCount:0, readerWait:0}},
			ClientIp:"",
			HttpRequestHeader:http.Header(nil)}
	*/

	// si := StringInterface()

	var srest []reflect.Value
	if s == "" {
		args := TestArgs{Name: "Jack", Value: "Okkk"}
		srest = c1.CallFunc(chan1, args)
	} else {
		srest = c1.CallFunc(chan1, s)
	}

	fmt.Printf("sret:%#v, length:%#v\n", srest, len(srest))
	if len(srest) > 0 {
		rr := srest[0] ///reflect.Value
		fmt.Printf("chan string:%#v \n", rr.Kind().String())
		inter := rr.Interface().(chan string)
		fmt.Printf("chan value:%#v \n", <-inter)
	}

	// reflect.Value 与 interface相互转换
	/*
		变量   --- reflect.ValueOf(i interface{}) --> reflect.Value
													 /
			\								 rV.interface()
			 类型断言							rV为 refelect.Value
				\								   /
					\							  /
						\						 /
								 interface{}
									类型
	*/
}

type TestArgs struct {
	Name  string
	Value string
}

func main() {
	//传入  "", 将使用 结构体 TestArgs的 Name字段
	CallFuncByReflect("")
}
