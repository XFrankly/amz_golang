package main

import "fmt"

//bool:                    %t
//int, int8 etc.:          %d
//uint, uint8 etc.:        %d, %#x if printed with %#v
//float32, complex64, etc: %g
//string:                  %s
//chan:                    %p
//pointer:                 %p
func main() {
	//	fmt.Printf("%d - %b - %x \n", 42, 42, 42)
	//	fmt.Printf("%d - %b - %#x \n", 42, 42, 42)
	//	fmt.Printf("%d - %b - %#X \n", 42, 42, 42)
	fmt.Printf("%d \t %b \t %#X \t %v \t %#v \n", 42, 42, 42, 42, 42)
}
