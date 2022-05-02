package main

import (
	"fmt"
	"log"

	"os"
)
//var log = logrus.New()
func init() {
	nf, err := os.Create("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	//log.Fatalln(nf)
	log.SetOutput(nf)  // 写日志到文件
}

func main() {
	_, err := os.Open("no-file.txt")
	if err != nil {
		//		fmt.Println("err happened", err)
		log.Println("err happened", err)
		//		log.Fatalln(err)
		//		panic(err)
	//
	//

	}
	f, err := os.OpenFile("testlogfile.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")


	//f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Println(err)
	//}
	//defer f.Close()
	//
	//logger := log.New(f, "prefix", log.LstdFlags)
	//logger.Println("text to append")
	//logger.Println("more text to append")
}

/*
Package log implements a simple logging package ... writes to standard error and prints the date and time of each logged message ... the Fatal functions call os.Exit(1) after writing the log message ... the Panic functions call panic after writing the log message.
*/

// Println calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
