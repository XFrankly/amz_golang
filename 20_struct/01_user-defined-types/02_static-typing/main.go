package main

import "fmt"

type foo int

func main() {
	var myAge foo //引用别名，不推荐
	myAge = 44
	fmt.Printf("%T %v \n", myAge, myAge)

	var yourAge int  // 推荐做法
	yourAge = 29
	fmt.Printf("%T %v \n", yourAge, yourAge)

	// this doesn't work:
	//	 fmt.Println(myAge + yourAge)

	// this conversion works:
	//	 fmt.Println(int(myAge) + yourAge)
}
