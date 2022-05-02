package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type person struct {
	first string
	last  string
	age   int
}

func main() {
	//结构体实例化
	p1 := person{"James", "Bond", 20}
	fmt.Println(p1)
	bs, _ := json.Marshal(p1)
	fmt.Printf("%T \n", bs)  //[]uint8  slice切片
	fmt.Println(reflect.TypeOf(bs)) //[]uint8  slice切片
	fmt.Println(string(bs))
}
