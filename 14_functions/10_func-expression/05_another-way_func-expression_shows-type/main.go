package main

import "fmt"
import "reflect"

func makeGreeter() func() string {
	return func() string {
		return "Hello world!"
	}
}

func main() {
	greet := makeGreeter()
	fmt.Println(greet())
	fmt.Printf("%T\n", greet) //show type
	fmt.Println(reflect.TypeOf(greet)) //show type

}
