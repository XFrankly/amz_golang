package main

import (
	"fmt"
	"gpc/envs"
)

var bindAddress = envs.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
	fmt.Println("hello")
	fmt.Printf("bind add str:%#v\n", bindAddress)
	envs.Parse()

	fmt.Printf("after parse:%#v\n", bindAddress)

}
