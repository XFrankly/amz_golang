package main

import "fmt"

type mySentences []string // 别名

func main() {
	var messages mySentences = []string{"Hello World!", "More coffee"}
	fmt.Println(messages)
	fmt.Printf("%T\n", messages)
	fmt.Printf("%T\n", []string(messages))
}
