package main

import "fmt"

func main() {
	var mp = map[string]int{
		"A": -1,
		"B": 0,
		"C": 1,
	}
	fmt.Println(mp["A"], mp["D"])

}
