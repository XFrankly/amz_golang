package main

import (
	"fmt"
	"reflect"
)

type Number interface {
	////泛型 数字 类型， 可以是 int 或 放loat
	int | float64
}

type GenericSlice[T any] []T /// 泛型接口类型，所有类型的切片

func MultiplyTen[T Number](a T) T {
	/// 接收所有数字类型并 返回
	return a * 10
}

type gener interface {
	int | float64
}
type Ordered interface {
	/*
		这个声明说Ordered接口是所有整数、浮点数和字符串类型的集合。竖线表示类型的联合（或本例中的类型集）。
		Integer并且是在包Float中类似定义的接口类型。constraints请注意，接口没有定义任何方法Ordered。
		类型约束中，通常不关心特定类型，例如string; 我们对所有字符串类型都感兴趣。这就是~令牌的用途
		~string 表示基础类型为所有类型集合string，这包括string本身 和所有使用定义声明的类型，例如
		type Mystring string

	*/
	int | float64 | ~string
}

///// 使用 已有第三方包
// func Max[E constraints.Ordered](ei []E) (max E) {
// 	// max := 0
// 	for _, i := range ei {
// 		if i > max {
// 			max = i
// 		}
// 	}
// 	return max
// }

type IOds struct {
	Id string
}

func (i *IOds) GetID() string {
	return i.Id
}

type Values struct {
	Id    string
	Value string
}

func (v *Values) GetID() string {
	return v.Id
}

type OrderTypes interface {
	IOds | Values
}

type GetIDs interface {
	GetID()
}

func MakeIOds() *IOds {
	return &IOds{
		Id: "test",
	}
}

//// 类型推断
////使用类型参数需要传递类型参数，这可能会产生冗长的代码。
///回到我们的通用GMin函数
func generals[T Ordered](a T) T {
	fmt.Printf("a:%T\n", a)
	return a
}

func generalsTwo[T Ordered](a T, b T) T {
	fmt.Printf("a:%T, b:%T\n", a, b)
	fmt.Println(a + b)
	return a + b
}

// func GetBySwitch[T GetIDs](t T) string {
// 	if &t == nil {
// 		return ""
// 	}
// 	return &t.GetID()
// }
func generalsType[T OrderTypes](a T) T {
	fmt.Printf("iods:%T \n", a)
	fmt.Printf("value id:%+v\n", a)
	// fmt.Println(a.Id)  // 虽然不能直接取值，可以接收 不同的类型
	if reflect.TypeOf(a).Name() == reflect.TypeOf(*MakeIOds()).Name() {
		// for r := range a {}
		fmt.Printf("%T,%T\n", a, &a)
		fmt.Println("typs string:", reflect.TypeOf(a).String())
		fmt.Println("Align:", reflect.TypeOf(a).Align())
		fmt.Println(reflect.TypeOf(a).FieldByName("Id")) // true 存在
		id_value := reflect.ValueOf(a).FieldByName("Id")
		fmt.Printf("%T,%+v \n", id_value, id_value) //reflect.Value,test  返回 Id 的值
		fmt.Println("receive:", &a, a)
	}
	var v2 interface{} = "大da君" // 将 string 类型赋值给 interface{}
	var v4 interface{} = v2
	fmt.Printf("%+v\n", v4)
	return a
}

func main() {
	a := 123
	var b float64
	b = 12.3
	generals(a)
	generals(b)
	/// 类型推断
	generalsTwo[float64](float64(a), b) // 全部 float
	generalsTwo[float64](1.21, 32.11)   // 显式调用
	generalsTwo[int](a, int(b))         // 全部 int

	// Tslice := make([]Ordered, 0)
	iod := MakeIOds()
	valud := &Values{
		Id:    "No1",
		Value: "Values",
	}
	generalsType(*iod)
	generalsType(*valud)

}
