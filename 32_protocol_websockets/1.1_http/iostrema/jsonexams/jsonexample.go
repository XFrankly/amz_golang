//package jsonexams
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
type Animal int

const (
	Unknown Animal = iota
	Gopher
	Zebra
)

func (a *Animal) UnmarshalJSON(b []byte) error {
	var sr string
	if err := json.Unmarshal(b, &sr); err != nil {
		return err
	}
	switch strings.ToLower(sr) {
	default:
		*a = Unknown
	case "gopher":
		*a = Gopher
	case "zebra":
		*a = Zebra
	}
	return nil
}

//func (a Animal) MarshalJSON() ([]byte, error) {
//	var sr string
//	switch a {
//	default:
//		sr = "unknown"
//	case Gopher:
//		sr = "gopher"
//	case Zebra:
//		sr = "zebra"
//	}
//	return json.Marshal(sr)
//}

func DePackageJsonList() {
	// 解析字典并 映射到 map[Animal]int
	blob := `["gopher", "armadillo", "zebra", "unknown", "gopher", "bee", "gopher", "zebra"]`
	var zoo []Animal
	if err := json.Unmarshal([]byte(blob), &zoo); err != nil {
		// 在发现 JSON 语法错误之前。 避免填写半个数据结构 检查格式是否正确。
		log.Fatalln(err)
	}
	census := make(map[Animal]int)
	for _, animal := range zoo {
		census[animal] += 1
	}

	fmt.Printf("Zoo Census: \n* Gopher: %d\n* Zebras: %d\n* Unknown: %d\n",
		census[Gopher], census[Zebra], census[Unknown])
	//控制台输出
	//Zoo Census:
	//* Gopher: 3
	//* Zebras: 2
	//* Unknown: 3
}

type Size int

const (
	Unrecognized Size = iota
	Small
	Large
)

func (s *Size) UnmarshalText(text []byte) error {
	// 不知道为什么 要绑定这个 函数给 Size ：？？？？？？？？？？？？？？
	switch strings.ToLower(string(text)) {
	default:
		*s = Unrecognized
	case "small":
		*s = Small
	case "large":
		*s = Large
	}
	return nil
}

func (s Size) MarshalText() ([]byte, error) {
	var name string
	switch s {
	default:
		name = "unrecognized"
	case Small:
		name = "small"
	case Large:
		name = "large"

	}
	return []byte(name), nil
}

func DePackageTextJson() {
	// 统计 列表 字符串 字符的出现次数
	blob := `["small","regular","large","unrecognized","small","normal","small","large"]`
	var inventory []Size
	if err := json.Unmarshal([]byte(blob), &inventory); err != nil {
		log.Fatal(err)
	}

	counts := make(map[Size]int)
	for _, size := range inventory {
		counts[size] += 1
	}

	fmt.Printf("Inventory Counts:\n* Small:        %d\n* Large:        %d\n* Unrecognized: %d\n",
		counts[Small], counts[Large], counts[Unrecognized])

}

///////////////////////////// Decoder  解析
func DecoderJsonStream() {
	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
	`

	type Message struct {
		// 将map 映射到 结构体
		Name, Text string
	}

	dec := json.NewDecoder(strings.NewReader(jsonStream))

	for {
		// 解析 map list
		var m Message
		if err := dec.Decode(&m); err == io.EOF { // 读完为止
			break
		} else if err != nil { // 解析错误，格式不匹配 结构体
			log.Fatalln(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)

	}

}

func DecoderDict() {
	// 解析 一层 dict信息
	//{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	const jsonResps = `
{"message":"Done", "token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre","status":200}
	`

	type TokenUser struct {
		token string
		user  string
	}

	type RespMap struct {
		// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}   //NoSupport
		// {"message":"Done", "token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre","status":200}
		Status  int
		Message string
		token   string
		user    string
	}

	decResp := json.NewDecoder(strings.NewReader(jsonResps))
	for {
		// 解析 map
		var Resp RespMap
		if err := decResp.Decode(&Resp); err == io.EOF {
			break
		} else if err != nil { // 解析错误，格式不匹配
			log.Println(err)
		}
		fmt.Printf("%v, status:%d \n", Resp.Message, Resp.Status)
	}
}

func MyDecoderMutilDict() {
	// 解析 多层 dict信息
	//{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	const jsonResps = `
	{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	`

	type RespMap struct {
		// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}   // Support
		Status  int
		Message map[string]string
	}
	fmt.Printf("jsonResps:%v", jsonResps)
	decResp := json.NewDecoder(strings.NewReader(jsonResps))
	for {
		// 解析 map
		var Resp RespMap
		if err := decResp.Decode(&Resp); err == io.EOF {
			break
		} else if err != nil { // 解析错误，格式不匹配
			log.Println(err)
		}
		fmt.Printf("%v, status:%d \n", Resp.Message, Resp.Status)
		fmt.Printf("token: %s, user:%s \n", Resp.Message["token"], Resp.Message["user"])
	}

}

// Token 解析，
func DecoderJsonToken() {
	// Token 混合解析, 不能解析中文
	// 按字符 逐个解析， 自动排除 关键字 { [
	// null, NUll 将被解析为 nil
	// 不能识别 python的 "None":none， None 将被识别为 Null
	const jsonStream = `
	{"Message": "Hello", "Array": [1, 2, 3], "Null": null, "None": "none", "Number": 1.234}
    `
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v", t, t)
		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}

}

func MyDecoderJsonToken() {
	// Token 混合解析, 不能解析中文
	// 按字符 逐个解析， 自动排除 关键字 { [
	// null, NUll 将被解析为 nil
	// 不能识别 python的 "None":none， None 将被识别为 Null
	const jsonStream = `
	{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
    `
	type RespMap struct {
		// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}   // Support
		Status  int
		Message []map[string]string
	}
	rm := &RespMap{}
	fmt.Println("respmap", rm)
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v", t, t)
		// 筛选并 填入 结构体实例

		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}

}

//HTML 标签编码
func CodeEscape() {
	var out bytes.Buffer
	json.HTMLEscape(&out, []byte(`{"Name":"<b>HTML content</b>"}`))
	out.WriteTo(os.Stdout)
}

// slice 切片 结构体 解码
func DecodeStructSlice() {
	type Road struct {
		Name   string
		Number int
	}
	roads := []Road{
		{"Diamond Fork", 29},
		{"Sheep Creek", 51},
	}
	b, err := json.Marshal(roads)
	if err != nil {
		log.Fatal(err)
	}
	var out bytes.Buffer
	json.Indent(&out, b, "=", " ") // 分隔符号  " " \n \t
	out.WriteTo(os.Stdout)
}

//DecodeMarshallIndent
func DecodeMarshall() string {
	// 结构体 转 为 混合类型 dict
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Printf("error:", err)
	}
	fmt.Printf("%T %v", b, b)
	os.Stdout.Write(b)
	return fmt.Sprintf("%s", b)
}

//DecodeMarshallIndent
func MyDecodeMarshall() string {
	// 结构体 转 为 混合类型 dict json
	// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	type UserToken struct {
		Token string
		User  string
	}
	type RespMap struct {
		// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}   // Support
		Status  int
		Message *UserToken //map[string]string
	}

	group := RespMap{
		Status:  1,
		Message: &UserToken{"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ==", "postgre"},
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Printf("error:", err)
	}
	fmt.Printf("%T %v", b, b)
	os.Stdout.Write(b)
	return fmt.Sprintf("%s", b)
}

func DecoderMarshalIndent() string {
	// 格式输出， <prefix>", "<indent>  + data
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	b, err := json.MarshalIndent(data, "<prefix>", "<indent>")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	return fmt.Sprintf("%s", string(b))
}

func DecodeRawMessageJson() string {
	// 使用 RawMessage, 为了在 封送期间  使用预先计算的 JSON
	//	 uses RawMessage to use a precomputed JSON during marshal.
	//{
	//        "header": {
	//                "precomputed": true
	//        },
	//        "body": "Hello Gophers!"
	//}
	//{
	//        "header": {
	//                "precomputed": true
	//        },
	//        "body": "Hello Gophers!"
	//}
	h := json.RawMessage(`{"precomputed":true}`)
	c := struct {
		Header *json.RawMessage `json:"header"`
		Body   string           `json:"body"`
	}{Header: &h, Body: "Hello Gophers!"}

	b, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	return string(b)
}

// unmarshal text
func DecodeByRawMsgUnmarshal() []interface{} {
	type Color struct {
		Space string
		Point json.RawMessage // 延迟解析 直到我们 得到 color 空间
	}
	type RGB struct {
		R uint8
		G uint8
		B uint8
	}
	type YCbCr struct {
		Y  uint8
		Cb int8
		Cr int8
	}
	var j = []byte(`[
	{"Space": "YCbCr", "Point": {"Y": 255, "Cb": 0, "Cr": -10}},
	{"Space": "RGB",   "Point": {"R": 98, "G": 218, "B": 255}}
]`)
	var colors []Color
	err := json.Unmarshal(j, &colors)
	if err != nil {
		log.Fatalln("error:", err)
	}

	var dst1 []interface{}
	for _, c := range colors {
		var dst interface{}
		switch c.Space {
		case "RGB":
			dst = new(RGB)
			dst1 = append(dst1, dst)
		case "YCbCr":
			dst = new(YCbCr)
			dst1 = append(dst1, dst)
		}
		err := json.Unmarshal(c.Point, dst)
		if err != nil {
			log.Fatalln("error:", err)
		}
		fmt.Println(c.Space, dst)
	}
	return dst1
}

///////////////////////////////
type Animals struct {
	Name  string
	Order string
}

var animals []Animals

func DecodeUnmarshall() []Animals {
	//func Unmarshal(data [] byte , v interface{}) error
	// 格式化输出
	//Unmarshal 解析 JSON 编码的数据并将结果存储在 v 指向的值中。如果 v 为 nil 或不是指针，则 Unmarshal 返回 InvalidUnmarshalError。
	//Unmarshal 使用 Marshal 使用的编码的逆编码，根据需要分配映射、切片和指针，并具有以下附加规则：
	//要将 JSON 解组为指针，
	//////Unmarshal 首先处理 JSON 为 JSON 文字 null 的情况。在这种情况下，Unmarshal 将指针设置为 nil。
	///////否则，Unmarshal 将 JSON 解组为指针指向的值。如果指针为 nil，Unmarshal 为其分配一个新值来指向。
	//要将 JSON 解组为实现 Unmarshaler 接口的值，Unmarshal 会调用该值的 UnmarshalJSON 方法，包括当输入为 JSON null 时
	//否则，如果该值实现 encoding.TextUnmarshaler 并且输入是 JSON 引用的字符串，则 Unmarshal 使用该字符串的未引用形式调用该值的 UnmarshalText 方法。
	//要将 JSON 解组到结构中，Unmarshal 将传入的对象键与 Marshal 使用的键（结构字段名称或其标记）匹配，首选完全匹配但也接受不区分大小写的匹配。
	//默认情况下，不具有相应结构字段的对象键将被忽略（请参阅 Decoder.DisallowUnknownFields 了解替代方案）。
	//为了将 JSON 解组为接口值，Unmarshal 将其中一项存储在接口值中：
	//	bool, 用于 JSON布尔值
	//float64，用于 JSON 数字
	//string，用于 JSON 字符串
	//[]interface{}，用于 JSON 数组
	//map[string]interface{}，用于 JSON 对象
	//nil JSON null

	//要将 JSON 数组解组为切片，
	//////Unmarshal 会将切片长度重置为零，然后将每个元素附加到切片。作为一种特殊情况，为了将空 JSON 数组解组为切片，Unmarshal 将切片替换为新的空切片。
	//要将 JSON 数组解组为 Go 数组，
	//////Unmarshal 将 JSON 数组元素解码为相应的 Go 数组元素。如果 Go 数组小于 JSON 数组，则丢弃额外的 JSON 数组元素。如果 JSON 数组小于 Go 数组，则额外的 Go 数组元素设置为零值。
	//要将 JSON 对象解组到映射中
	/////Unmarshal 首先建立要使用的映射。如果映射为 nil，Unmarshal 分配一个新映射。否则 Unmarshal 会重用现有映射，保留现有条目。然后 Unmarshal 将 JSON 对象中的键值对存储到映射中。
	////映射的键类型必须是任何字符串类型、整数、实现 json.Unmarshaler 或实现 encoding.TextUnmarshaler。
	//如果 JSON 值不适合给定的目标类型，或者 JSON 编号溢出目标类型，Unmarshal 会跳过该字段并尽其所能完成解组。
	////如果没有遇到更严重的错误，Unmarshal 返回一个 UnmarshalTypeError 描述最早的此类错误。在任何情况下，都不能保证有问题的字段后面的所有剩余字段都将被解组到目标对象中。
	//JSON null 值通过将 Go 值设置为 nil 来解组为接口、映射、指针或切片。因为 null 在 JSON 中经常用来表示“不存在”，所以将 JSON null 解组到任何其他 Go 类型对值没有影响并且不会产生错误。
	//解组带引号的字符串时，无效的 UTF-8 或无效的 UTF-16 代理对不会被视为错误。相反，它们被 Unicode 替换字符 U+FFFD 替换。
	var jsonBlob = []byte(`[
	{"Name":"Platypus", "Order":"Monotremata"},
		{"Name":"Quoll", "Order":"Dasyuromorphia"}
	]
`)

	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	return animals
}

/// 检查 Json 格式 是否 正确
func CheckJsonFormat(jsonstr string) (bool, bool) {
	goodJson := `{"example":1}`
	badJson := `{"example":2[}`

	v1 := json.Valid([]byte(goodJson))
	v2 := json.Valid([]byte(badJson))
	fmt.Println()
	fmt.Println()
	return v1, v2
}
func main() {
	// 解析package list
	//DePackageJsonList()

	// 解析 package Text
	//DePackageTextJson()

	//  解析 多个 map 流
	//DecoderJsonStream()

	// 解析一层map
	//DecoderDict()

	//把 多层结构化的数据 解析为 map
	MyDecoderMutilDict()

	// 逐个字符解析
	//DecoderJsonToken()
	// Token 解析混合类型
	//MyDecoderJsonToken()

	////HTML 标签编码
	//CodeEscape()

	// slice 切片 结构体 解码
	//DecodeStructSlice()

	// 结构体 转 为 混合类型 dict  示例
	//print("\n", DecodeMarshall())

	//结构体转为混合类型 字典 map
	//print("\n", MyDecodeMarshall())

	//// 格式输出， <prefix>", "<indent>  + data
	//print("\n", DecoderMarshalIndent())

	//print("\n", DecodeRawMessageJson())
	//print("unmarshal : \n", fmt.Sprintf("%s", DecodeByRawMsgUnmarshal()))

	//print("\nDecodeUnmarshall:\n", fmt.Sprintf("%s", DecodeUnmarshall()))
	// 判定 JSON 编码格式是否正确
	tf, tf2 := CheckJsonFormat("")
	print(fmt.Sprintf("%T, %s\n", tf, tf))
	fmt.Sprintf("%T, %s", tf2)
}
