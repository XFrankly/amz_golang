package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

/*
Call() 用于调用方法，传入方法的 参数
		注意区别，go-rpc 自带框架实现
		reflect.ValueOf().Call() 和 reflect.Method().Func.Call()
		这两个call不同的是，ValueOf() 不需要传递receiver
		而Method() 第一个参数必须是 接收器 receiver

		示例，通过反射调用函数
			//将函数包装为反射值对象
		    funcValue := reflect.ValueOf(add)

		    //构造函数参数，传入两个整形值
		    paramList := []reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)}

		    //反射调用函数
		    retList := funcValue.Call(paramList)

		    fmt.Println(retList[0].Int())
*/
type T struct{}

func (t *T) Infos() bool {
	fmt.Println("Call infos.")
	return true
}

func CallExample1() {
	var t T
	//使用reflect Call调用
	val := reflect.ValueOf(&t).MethodByName("Infos").Call([]reflect.Value{})
	fmt.Println(val)
}

func (t *T) Infos2() string {
	var msg string
	stdMsg := "Golang Infos2."
	msg += stdMsg
	fmt.Println(msg)
	stdMsg2 := "reflect. Call Function"
	msg += stdMsg2
	fmt.Println(msg)
	return msg
}

func CallEx2() {
	const num = 2
	var t T
	var t1 T
	tval := reflect.ValueOf(&t).MethodByName("Infos").Call([]reflect.Value{})
	t1val := reflect.ValueOf(&t1).MethodByName("Infos2").Call([]reflect.Value{})
	fmt.Println("tval:", tval, "t1val:", t1val)
	fmt.Printf("tval:%#v, val 0 string:%#v\n", tval[0], tval[0].String())
}

/////////////////////使用 reflect.Call 调用全部类型的某个方法
type Base struct{}
type Teacher struct {
	Base
	Name      string
	BirthYear int
	Grander   string
	Item      string
}

func NewTeacher(n, g, i string, b int) *Teacher {
	if n == "" {
		n = "Jack"
	}

	if g == "" {
		g = "Class5"
	}

	if i == "" {
		i = "Math"
	}
	if b == 0 {
		b = 45
	}
	return &Teacher{
		Name:      n,
		BirthYear: b,
		Grander:   g,
		Item:      i,
	}
}
func (t *Teacher) BYear() int {
	return t.BirthYear
}
func (t *Teacher) Age() int {
	c := time.Now().Year()
	age := c - t.BYear()
	return age
}
func (t *Teacher) IsActive(force int) bool {
	var live bool
	if force == 0 {
		live = true
	}
	if t.Grander == "" && t.Item == "" {
		live = false
	}
	live = true
	fmt.Printf("teacher %v is work? %v \n", t.Name, live)
	return live
}

type School struct {
	Base
	Name      string //名称
	BuildYear int    //开办年份
	Type      string //学校类型
	Students  int    //总人数
}

func (s *School) Age() int {
	return time.Now().Year() - s.BuildYear
}

func (s *School) IsActive(force int) bool {
	var live bool
	if force == 0 {
		live = true
	}
	if s.Students == 0 {
		live = false
	}
	live = true
	fmt.Printf("School %v is work? %v \n", s.Name, live)

	return live
}

func CaseCallExample3(t interface{}) any {
	//Age 方法不需要传入参数值
	tval := reflect.ValueOf(t).MethodByName("Age").Call([]reflect.Value{})
	fmt.Printf("age tval:%#v\n", tval[0])
	force := 0
	refVal := reflect.ValueOf(force)
	val := []reflect.Value{}
	val = append(val, refVal)
	//IsActive 方法需要传入 参数 force int 类型
	aliveVal := reflect.ValueOf(t).MethodByName("IsActive").Call(val)
	fmt.Println("is Alive?:", aliveVal[0])
	return tval[0]
}

// 使用 接口实 现 对多种类型的 相同方法的调用
type CallInter interface {
	Age() int
	IsActive(force int) bool
	DoCall(ci CallInter) bool
}

func (b *Base) DoCall(ci CallInter) bool {
	a := ci.Age()

	alive := ci.IsActive(0)
	fmt.Println("its type:", reflect.TypeOf(ci).String(), "age:", a, "is alive?:", alive)
	return alive
}

func main() {
	CallExample1()
	CallEx2()

	t1 := NewTeacher("", "", "", 1945)
	CaseCallExample3(t1)
	s1 := &School{Name: "BeiKing Universon", BuildYear: 1949, Type: "All", Students: 19999}
	CaseCallExample3(s1)

	//使用接口
	var c = make(chan bool, 2)
	var wg sync.WaitGroup
	n := &Base{}
	wg.Add(2)
	go func() {

		r := n.DoCall(t1)
		c <- r
		wg.Done()
	}()
	go func() {
		r := n.DoCall(s1)
		c <- r
		wg.Done()
	}()
	wg.Wait()

	// 从通道取值
	for len(c) > 0 {
		select {
		case c1 := <-c:
			fmt.Println("chan result:", c1)

		}
	}

}
