package main

/*
结构体到json的映射
*/

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/// {"name":"bob", "age":10}
type User struct {
	Name string `json:"name" zh:"姓名"`
	Age  int64  `json:"age" zh:"年龄"`
}
type City struct {
	Name       string `en:"xname" json:"name"` //// 在解析为json 格式时，将 使用该名字 json:name
	Population int64  `en:"xpopulation", json:"population"`
	GDP        int64  `en:"xgdp", json:"gdp"`
	Mayor      string `en:"xmayor" json:"mayor"`
	Active     bool   `en:"xactive" json:"active"`
}

func reflectExample() {
	var x float64 = 3.14
	var u User = User{"bod", 10}

	fmt.Println(x, u)

	refObjValPtr := reflect.ValueOf(&x) /// 指针的方式更准确
	refObjType := reflect.TypeOf(&x)
	refObVal := reflect.ValueOf(x)
	refObTyp := reflect.TypeOf(x)
	// refObTyp2 := reflect.(x)
	// fmt.Printf("refObjValPrt float64 x:%f\n", refObjValPtr.Float()) // 不能直接取指针的值，需要先使用elem
	fmt.Printf("reflect type:%+v\n", refObjType.String()) //*float64
	fmt.Printf("ref type kind:%v\n", refObjType.Kind())   //ptr

	refObjVal := refObjValPtr.Elem() /// 返回指针所在位置的 变量

	fmt.Printf("ref value for  x:%s\n", refObjVal.String())
	fmt.Printf("ref type for x:%s\n", refObTyp.String())
	fmt.Printf("value float64 x:%f\n", refObVal.Float())

	fmt.Printf("value2 float64 x:%v\n", refObjVal.Float())
	fmt.Printf("value canset x:%v\n", refObjVal.CanSet())

	refObjVal.Set(reflect.ValueOf(4.25)) // 更新值
	fmt.Printf("updated x using point*refObjVal %+v\n", refObjVal)

	urefObVal := reflect.ValueOf(u)
	urefObTyp := reflect.TypeOf(u)
	fmt.Println("ref value obj for u:", urefObVal)
	fmt.Println("ref type obj for u:", urefObTyp)
	fmt.Println("ref kind obj for u:", urefObTyp.Kind(), urefObVal.Kind())
}

func JSONEncode(v interface{}, tagKey string) ([]byte, error) {
	refObjval := reflect.ValueOf(v)

	refObjType := reflect.TypeOf(v)
	buf := bytes.Buffer{}
	if refObjval.Kind() != reflect.Struct {
		return buf.Bytes(), fmt.Errorf(
			"val of kind %s not support",
			refObjval.Kind(),
		)
	}
	buf.WriteString("{")
	pairs := []string{}
	for i := 0; i < refObjval.NumField(); i++ {
		structFieldRefObj := refObjval.Field(i)
		structFieldRefObjTyp := refObjType.Field(i)

		tag := structFieldRefObjTyp.Tag.Get(tagKey)
		switch structFieldRefObj.Kind() {
		case reflect.String:
			//
			strVal := structFieldRefObj.Interface().(string)
			// pairs = append(pairs, `"`+string(structFieldRefObjTyp.Tag)+`":`+strVal)
			pairs = append(pairs, `"`+tag+`":`+strVal)
		case reflect.Int64:
			//
			intVal := structFieldRefObj.Interface().(int64)
			// pairs = append(pairs, `"`+string(structFieldRefObjTyp.Tag)+`":`+strconv.FormatInt(intVal, 10))
			pairs = append(pairs, `"`+tag+`":`+strconv.FormatInt(intVal, 10))
		case reflect.Bool:
			boolVal := structFieldRefObj.Interface().(bool)
			// pairs = append(pairs, `"`+string(structFieldRefObjTyp.Tag)+`":`+strconv.FormatBool(boolVal))
			pairs = append(pairs, `"`+tag+`":`+strconv.FormatBool(boolVal))

		default:
			//
			return buf.Bytes(), fmt.Errorf(
				"struct field with name %s and kind %s is not support",
				structFieldRefObjTyp.Name, // 属性 非函数
				structFieldRefObj.Kind(),
			)
		}
	}
	// pairs = []string{""}
	buf.WriteString(strings.Join(pairs, ","))
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func UseJsonFormat() {
	var u User = User{"Pop", 20}
	c := City{Name: "Sf", Population: 3000000, GDP: 1000000, Mayor: "joe", Active: true}
	res, err := JSONEncode(u, "zh") /// 取json 标签的名称
	fmt.Println("error:", err)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	cres, err1 := JSONEncode(c, "en") /// 有多个标签 取 en 标签的名称
	fmt.Println("error:", err1)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(cres))
}
func main() {
	// reflectExample()
	UseJsonFormat()

}
