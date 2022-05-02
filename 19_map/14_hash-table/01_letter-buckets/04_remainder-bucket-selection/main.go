package main

import "fmt"

func main() {
	fmt.Println(string(0))
	for i := 65; i <= 122; i++ {
		fmt.Println(i, " - ", string(i), " - ", i%12)
	}
}
