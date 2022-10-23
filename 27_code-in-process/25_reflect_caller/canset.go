package main

import (
	"fmt"
	"reflect"
)

/*
### 值可修改条件之一：可被寻址
		通过反射修改变量值的前提条件之一：这个值必须可以被寻址。简单地说就是这个变量必须能被修改。示例代码如下
		 //声明整形变量a并赋初值
		 1   var a int = 1024

		    //获取变量a的反射值对象
		 2    rValue := reflect.ValueOf(a)

		    //尝试将a修改为1(此处会崩溃)
		 3   rValue.SetInt(1)

		 由于2 传入的是a的值，无法获取在 第三步 设置 a的新值

		 只需要修改 2 为a的地址，即指针即可
		  2    rValue := reflect.ValueOf(&a)

	### 值可修改的条件二：被导出，即变量首字母大写
		结构体成员中，如果字段没有被导出，即便不使用反射也可能被访问，但不能被修改
		type dog struct {
      		  LegCount int
      		  littleCount int
   			 }

		    //获取dog实例的反射值对象
		    valueOfDog := reflect.ValueOf(&dog{})

		  //// 取出dog实例地址的元素
		    valueOfDog = valueOfDog.Elem()

		    //获取legCount字段的值,
		    vLegCount := valueOfDog.FieldByName("LegCount")
 			//获取littleCount字段的值,
		    vlittleCount := valueOfDog.FieldByName("littleCount")

		    //尝试设置legCount的值 由于LegCount首字母大写，已导出可以设置成功
		    vLegCount.SetInt(4)

		    //尝试设置vlittleCount 值失败，由于 littleCount 首字母小写，未导出
		    vlittleCount.SetInt(1)

		    fmt.Println(vLegCount.Int())

*/
//Canset
type ProductionInfo struct {
	StructA []Entry
}

type Entry struct {
	Field1 string
	Field2 int
}

func SetField(source interface{}, fieldName string, fieldValue interface{}) {
	v := reflect.ValueOf(source)
	tt := reflect.TypeOf(source)

	for k := 0; k < tt.NumField(); k++ {
		fieldValue := reflect.ValueOf(v.Field(k))

		// use of CanSet() method
		fmt.Println(fieldValue.CanSet())
		if fieldValue.CanSet() {
			fieldValue.SetString(fieldValue.String())
		}
	}
}

func SetField2(source interface{}, fieldName string, fieldValue string) {
	v := reflect.ValueOf(source).Elem()

	// use of CanSet() method
	fmt.Println(v.FieldByName(fieldName).CanSet())

	if v.FieldByName(fieldName).CanSet() {
		v.FieldByName(fieldName).SetString(fieldValue)
	}
}

/*
reflect CanSet
報告是否可以更改 v 的值。
// 只有當它是可尋址的並且不可尋址時，才能改變它
// 通過使用未導出的結構字段獲得。
// 如果 CanSet 返回 false，則調用 Set 或任何特定類型
// setter (e.g., SetBool, SetInt) 會恐慌。
*/
type Poeple struct {
	Name string
	Age  int
}

func CaseCansetTrue() {
	source := ProductionInfo{}
	p := Poeple{"Jack", 44}
	source.StructA = append(source.StructA, Entry{Field1: "A", Field2: 2})

	fmt.Println("Before: ", source.StructA[0])
	SetField2(&source.StructA[0], "Field1", "NEW_VALUE")
	fmt.Println("After: ", source.StructA[0])

	fmt.Println("Before People Name:", p)
	SetField2(&p, "Name", "Lucy")
	fmt.Println("After People Name:", p)
	p.Name = "Dave"
	fmt.Println("End:", p)
}
func main() {

	source := ProductionInfo{}
	source.StructA = append(source.StructA, Entry{Field1: "A", Field2: 2})

	SetField(source, "Field1", "NEW_VALUE")
	SetField(source, "Total", 2)

	CaseCansetTrue()
}
