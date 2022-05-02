package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	rdr := strings.NewReader("test")
	fmt.Println("rdr", rdr)
	io.Copy(os.Stdout, rdr)
	fmt.Println(os.Stdout)
}
