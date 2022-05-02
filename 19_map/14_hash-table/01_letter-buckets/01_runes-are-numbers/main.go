package main

import "fmt"

func main() {
	letter := 'A'
	letter2 := "A"
	fmt.Println(letter, letter2)
	ms := map[int]string{
		1:"one",
		2:"two",
	}
	fmt.Printf("%T %T \n", letter, letter2)
	fmt.Printf("%T", ms)
}
