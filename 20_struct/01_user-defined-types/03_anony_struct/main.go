package main

import "fmt"

type student struct {
	name string
	age  int
}

//嵌套
func selfNested() {
	m := make(map[string]*student)
	stus := []student{
		{name: "小王子", age: 18},
		{name: "娜扎", age: 23},
		{name: "大王", age: 9000},
	}

	for _, stu := range stus {
		stu := stu //如果没有这句，将导致所有stu都为最后一个大王
		m[stu.name] = &stu
	}
	for k, v := range m {
		fmt.Println(k, "=>", v.name)
	}
}

//结构体匿名字段
//Address 地址结构体
type Address struct {
	Province   string
	City       string
	CreateTime string
}

//Email 邮箱结构体
type Email struct {
	Account    string
	CreateTime string
}

//User 用户结构体
type User struct {
	Name   string
	Gender string
	Address
	Email
}

///匿名字段冲突时，需要显式指定
func mainAns() {
	var user3 User
	user3.Name = "沙河娜扎"
	user3.Gender = "男"
	// user3.CreateTime = "2019" //ambiguous selector user3.CreateTime
	user3.Address.CreateTime = "2000" //指定Address结构体中的CreateTime
	user3.Email.CreateTime = "2000"   //指定Email结构体中的CreateTime
	fmt.Printf("user3:%v\n", user3)   //user3:{沙河娜扎 男 {  2000} { 2000}}
}

//匿名结构体
func main() {
	var user struct {
		Name string
		Age  int
	}
	user.Name = "Jack"
	user.Age = 22
	fmt.Printf("%v\n", user)

	// 空结构体
	var em struct{}
	fmt.Printf("%T, %#v, %v\n", em, em, &em)

	selfNested()
	/*
		娜扎 => 娜扎
		大王 => 大王
		小王子 => 小王子
	*/

	mainAns()
}
