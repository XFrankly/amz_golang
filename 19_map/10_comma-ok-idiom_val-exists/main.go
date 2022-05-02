package main

import (
	"fmt"
)

func main() {

	myGreeting := map[int]string{
		0: "Good morning!",
		1: "Bonjour!",
		2: "Buenos dias!",
		3: "Bongiorno!",
	}

	fmt.Println(myGreeting)

	// delete(myGreeting, 2)
	myGreeting[8] = ""
	if val, exists := myGreeting[9]; exists {
		fmt.Println("if That value exists.")
		fmt.Println("if val: ", val)
		fmt.Println("if exists: ", exists)
	} else {
		fmt.Println("else That value doesn't exist.")
		fmt.Println("else val: ", val)
		fmt.Println("else exists: ", exists)
	}

	fmt.Println(myGreeting)
}
