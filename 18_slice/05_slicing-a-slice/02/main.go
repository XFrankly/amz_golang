package main

import "fmt"
import "reflect"

func main() {

	greeting := []string{
		"Good morning!",
		"Bonjour!",
		"dias!",
		"Bongiorno!",
		"Ohayo!",
		"Selamat pagi!",
		"Gutten morgen!",
	}

	fmt.Print("[1:2] ")
	fmt.Println(greeting[1:2])
	fmt.Print("[:2] ")
	fmt.Println(greeting[:2])
	fmt.Print("[5:] ")
	fmt.Println(greeting[5:])
	fmt.Print("[:] ")
	fmt.Println(greeting[:])
	greeting[0] = "12312"
	fmt.Println(greeting[:])

	for i, t := range greeting {
		fmt.Println(i, t)
		//var t1 int64 = (int64)(unsafe.Pointer(t)
		fmt.Println(reflect.TypeOf(t), []byte(t)) //, t1)
	}
}
