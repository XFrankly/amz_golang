package main

import (
	"io"
	"os"
)

//只读src.txt的前5个字符
func main() {
	src, err := os.Open("src.txt")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create("dst.txt")
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	rdr := io.LimitReader(src, 5)
	io.Copy(dst, rdr)

}
