package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 从已有的文件 复制一个 相同的
func main() {
	f, err := os.Open(os.Args[1])
	fmt.Println("f", f)
	if err != nil {
		log.Fatalln("my program broke opening: ", err.Error())
	}
	defer f.Close()

	nf, err := os.Create("newFile.txt")
	if err != nil {
		log.Fatalln("my program broke creating: ", err.Error())
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("my program broke reading: ", err.Error())
	}

	_, err = nf.Write(bs)
	if err != nil {
		log.Fatalln("my program broke writing: ", err.Error())
	}
}

/*

step 1 - at command line enter:
go install

step 2 - at command line enter:
07_copy main.go

*/
