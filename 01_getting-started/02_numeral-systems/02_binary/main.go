package main

import "fmt"
//bool:                    %t
//int, int8 etc.:          %d
//uint, uint8 etc.:        %d, %#x if printed with %#v
//float32, complex64, etc: %g
//string:                  %s
//chan:                    %p
//pointer:                 %p
//binary                   %b
func main() {
	fmt.Printf("%d - %b \n", 42, 42)
	fmt.Printf("%d -  \n",52)
}
