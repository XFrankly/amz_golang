package main

import (
	"fmt"
	"reflect"
)

func MakeFuncsTests() {
	// swap  傳遞給 MakeFunc 的實現。
	// 它必鬚根據 reflect.Values 工作，這樣才有可能
	// 將會在事先不知道類型的情況下編寫代碼
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}

	// makeSwap 期望 fptr 是指向 nil 函數的指針。
	// 它將指針設置為使用 MakeFunc 創建的新函數。
	// 當函數被調用時，reflect 會轉換參數
	// 轉化為Values，調用swap，然後將swap的結果切片
	// 到新函數返回的值中。
	makeSwap := func(fptr any) {
		// fptr 是指向函數的指針。
		// 以 reflect.Value 的形式獲取函數值本身（可能是 nil）
		// 這樣我們就可以查詢它的類型，然後設置值。
		fn := reflect.ValueOf(fptr).Elem()

		// 製作正確類型的函數。
		v := reflect.MakeFunc(fn.Type(), swap)

		// 將其賦值給 fn 表示的值。
		fn.Set(v)
	}

	// 為整數創建並調用交換函數。
	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	fmt.Println(intSwap(0, 1))

	//為 float64s 創建並調用交換函數。
	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(2.72, 3.14))

}

type TestArg struct {
	Name  string
	Value string
}

//大类型断定
func KindReflect(args interface{}) {
	var val reflect.Value

	//如果是golang 数据类型，直接返回值，否则需要使用 Elem()
	fmt.Printf("args all elem:%#v\n", reflect.ValueOf(args))

	if reflect.TypeOf(args).Name() == reflect.TypeOf(TestArg{}).Name() {
		//Args结构体
		targs := args.(TestArg)
		val = reflect.ValueOf(targs.Name)
	} else if reflect.TypeOf(args).Name() == reflect.TypeOf("a").Name() {
		//字符
		val = reflect.ValueOf(args.(string))
	} else if reflect.TypeOf(args).Kind() == reflect.TypeOf(struct{}{}).Kind() {
		//匿名结构体
		fmt.Println("its a struct{}{}", args)
	} else if reflect.TypeOf(args).Kind() == reflect.TypeOf(&struct{}{}).Kind() {
		//匿名指针
		fmt.Println("its a ptr &struct{}{}", args)

	} else {
		//未知场景

		val = reflect.ValueOf("Hello, Whats That?")
		fmt.Println("will panic:", val)
		panic(args)

	}
	fmt.Println("Call with args:", val)
}
func main() {

	KindReflect(TestArg{Name: "Args", Value: "111"})  //args
	KindReflect("Haha")                               //string
	KindReflect(struct{}{})                           //struct
	KindReflect(&struct{}{})                          //ptr
	KindReflect(&TestArg{Name: "Args", Value: "111"}) //ptr

}
