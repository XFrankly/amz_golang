package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("os args:", os.Args)
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("my program broke: ", err.Error())
	}
	defer f.Close()

	io.Copy(os.Stdout, f)
}
