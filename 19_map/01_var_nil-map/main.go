package main

import "fmt"

func main() {

	var myGreeting map[string]string  // 只是声明，没有赋值
	myGreeting = myGreeting
	//myGreeting["Tim"] = "Good morning."
	fmt.Println(myGreeting)
	fmt.Println(myGreeting == nil)

	var myGreetingi  = make(map[int]int)
	myGreetingi[0] = 1
	myGreetingi[0] = 1
	myGreetingi[1] = 2
	myGreetingi[3] = 3


	fmt.Println(myGreetingi)
}

// add these lines:
/*
	myGreeting["Tim"] = "Good morning."
	myGreeting["Jenny"] = "Bonjour."
*/
// and you will get this:
// panic: assignment to entry in nil map
