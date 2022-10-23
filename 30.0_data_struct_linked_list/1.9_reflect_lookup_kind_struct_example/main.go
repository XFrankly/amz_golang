package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

/*
Kind 表示 Type 表示的特定类型的类型。零类型不是有效类型。
*/

func kindExample() {
	for _, v := range []any{"hi", 42, func() {}} {
		switch v := reflect.ValueOf(v); v.Kind() {
		case reflect.String:
			fmt.Println(v.String())
		case reflect.Float32, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Println(v.Int())
		default:
			fmt.Printf("unhandled kind:%s", v.Kind())
		}
	}
}

/*
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value

MakeFunc 返回一个包含函数 fn 的给定类型的新函数。调用时，该新函数执行以下操作：

- 将其参数转换为值的 一部分。
- 运行结果：= fn(args)。
- 将结果作为值的一部分返回，每个正式结果一个。

实现 fn 可以假设参数 Value 切片具有由 typ 给出的参数数量和类型。
如果 typ 描述了一个可变参数函数，那么最终的 Value 本身就是一个表示可变参数参数的切片，
就像在可变参数函数的主体中一样。
fn 返回的结果值切片必须具有 typ 给定的结果数量和类型。

Value.Call 方法允许调用者根据值调用类型化函数；
相反，MakeFunc 允许调用者根据值实现类型化函数。

文档的示例部分包含如何使用 MakeFunc 为不同类型构建交换函数的说明。
*/

func MakeFuncExample() {
	// swap 是传递给 MakeFunc 的实现
	// 如果可以 它必须在 工作在 reflect.Values的 语句中，
	// 以便提供一个可能： 在不知道 类型的场景 写code
	// swap 是传递给 MakeFunc 的实现。
	// 它必须根据 reflect.Values 工作，这样才有可能
	// 在事先不知道类型的情况下编写代码
	// 将会。
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}

	/*
		 makeSwap 期望 fptr 是一个指向 nil 函数的指针。
		它将指针设置为使用 MakeFunc 创建的新函数。
			当函数被调用时，reflect 会转换参数转化为Values，调用swap，
			然后将swap的结果切片
			到新函数返回的值中
	*/
	makeSwap := func(fptr any) {
		// fptr 是一个指向函数的指针
		// // fptr 是一个指向函数的指针。
		// 以 reflect.Value 的形式获取函数值本身（可能是 nil）
		// 这样我们就可以查询它的类型，然后设置值。
		fn := reflect.ValueOf(fptr).Elem()

		// 创建一个正确类型的 函数
		v := reflect.MakeFunc(fn.Type(), swap)

		// 将其赋值给 fn 代表的值
		fn.Set(v)
	}

	// 为 ints 创建和调用一个 swap 函数
	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	fmt.Println(intSwap(2, 3)) /// 交换顺序

	//为float64 创建一个 swap 函数
	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(2.72, 3.14)) /// 交换顺序
}

/*
func StructOf(fields []StructField) Type
StructOf 返回包含字段的结构类型。 Offset 和 Index 字段将被编译器忽略和计算。
如果传递了未导出的 StructField，StructOf 当前不会为嵌入字段生成包装器方法和恐慌。
这些限制可能会在未来的版本中取消。
*/

func StructOfExample() {
	typ := reflect.StructOf([]reflect.StructField{
		{Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height`},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age`,
		},
	})

	v := reflect.New(typ).Elem()
	v.Field(0).SetFloat(0.5)
	v.Field(1).SetInt(2)
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	fmt.Printf("value:%+v\n", s)
	fmt.Printf("json:%s\n", w.Bytes())

	r := bytes.NewReader([]byte(`{"height":1.7,"age":18}`))
	if err := json.NewDecoder(r).Decode(s); err != nil {
		panic(err)
	}
	fmt.Printf("value:%+v\n", s)
}

//// TypeOf 返回表示 i 的动态类型的反射 Type。如果 i 是一个 nil 接口值，TypeOf 返回 nil。

func TypeOfExample() {
	// 类似一个 接口类型，仅仅被静态类型使用
	/// 一个查找接口的反射类型的常用习语
	// Foo 类型是使用 *Foo 值
	writerType := reflect.TypeOf((*io.Writer)(nil)).Elem()

	fileType := reflect.TypeOf((*os.File)(nil))
	/// Implements 报告该类型是否实现了接口类型 u。
	fmt.Println(fileType.Implements(writerType))
}

////StructTag StructTag 是结构字段中的标记字符串。
//按照惯例，标签字符串是可选用空格分隔的键：“值”对的串联。每个键都是一个非空字符串，
//由除空格 (U+0020 ' ')、引号 (U+0022 '"') 和冒号 (U+003A ':') 以外的非控制字符组成。
//每个值都被引用使用 U+0022 '"' 字符和 Go 字符串文字语法。
func TypeOfStructTag() {
	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
}

////
func FieldBuIndexExample() {
	/*
			FieldByIndex 返回索引对应的嵌套字段。如果评估需要单步执行 nil 指针或不是结构的字段，它会出现错误。
		// 这个例子展示了一个提升字段名称的情况
		// 被另一个字段隐藏：FieldByName 不起作用，所以
		// 必须改为使用 FieldByIndex。
	*/
	type user struct {
		firstName string
		lastName  string
	}

	type data struct {
		user
		firstName string
		lastName  string
	}

	u := data{
		user:      user{"Embeded johb", "Embededd doe"},
		firstName: "john",
		lastName:  "Doe",
	}

	s := reflect.ValueOf(u).FieldByIndex([]int{0, 1})
	fmt.Println("embeded last name:", s)
}

//// reflect.lookup
/*
查找返回与标签字符串中的键关联的值。如果标签中存在键，则返回值（可能为空）。
否则返回值将是空字符串。 ok 返回值报告该值是否在标记字符串中显式设置。
如果标签没有常规格式，则 Lookup 返回的值是未指定的。
*/
func LookupExample() {
	type S struct {
		F0 string `alias:"field_0"`
		F1 string `alias:""`
		F2 string
	}
	s := S{}
	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if alias, ok := field.Tag.Lookup("alias"); ok {
			if alias == "" {
				fmt.Println("(blank)")
			} else {
				fmt.Println(alias)
			}
		} else {
			fmt.Println("(not specified)")
		}
	}
}

func main() {
	kindExample()

	//为不同类型构建交换函数
	MakeFuncExample()

	/// StructOf 返回包含字段的结构类型。 Offset 和 Index 字段将被编译器忽略和计算。
	StructOfExample()

	// 类型检查 示例
	TypeOfExample()
	/// 按索引取结构体的值
	FieldBuIndexExample()
	// 结构体 属性标签值的 获取
	TypeOfStructTag()
	/// 查找返回与标签字符串中的键关联的值
	LookupExample()
}
