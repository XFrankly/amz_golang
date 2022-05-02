package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type person struct {
	First string
	Last  string
	Age   int `json:"wisdom score"`
}

func main() {
	var p1 person
	fmt.Println(p1.First)
	fmt.Println(p1.Last)
	fmt.Println(p1.Age)
	//定义 分片
	bs := []byte(`{"First":"James", "Last":"Bond", "wisdom score":20}`)
	fmt.Println(reflect.TypeOf(bs))  //[]uint8
	fmt.Printf("%T \n", bs)
	json.Unmarshal(bs, &p1) //反序列化

	fmt.Println("--------------")
	fmt.Println(p1.First)
	fmt.Println(p1.Last)
	fmt.Println(p1.Age)
	fmt.Printf("%T \n", p1)
}
