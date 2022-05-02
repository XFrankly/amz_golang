package main

import "fmt"

func main() {

	myGreeting := map[string]string{
		"zero":  "Good morning!",
		"one":   "Bonjour!",
		"two":   "Buenos dias!",
		"three": "Bongiorno!",
	}
	myGreeting["8"] = ""
	fmt.Println(myGreeting)
	fmt.Printf("myGreeting format : %T", myGreeting)
	delete(myGreeting, "two")
	fmt.Println(myGreeting)
}
