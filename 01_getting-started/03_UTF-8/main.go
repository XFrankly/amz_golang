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
	for i := 60; i < 122; i++ {
		fmt.Printf("%d \t %b \t %x \t %q \n", i, i, i, i)
	}
}
