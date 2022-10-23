package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {

	from := "fram.gaya@gmail.com"
	to := "kuke.gaia@gmail.com"
	name := "John gaya"
	subject := "fra hello"
	body := "For job title."

	host := "core9:25" /// 邮件服务器地址

	con, err := net.Dial("tcp", host)
	checkError(err)

	req := "HELO core9\r\n" +
		"MAIL FROM: " + from + "\r\n" +
		"RCPT TO: " + to + "\r\n" +
		"DATA\r\n" +
		"From: " + name + "\r\n" +
		"Subject: " + subject + "\r\n" +
		body + "\r\n.\r\n" + "QUIT\r\n"

	_, err = con.Write([]byte(req))
	checkError(err)

	res, err := ioutil.ReadAll(con)
	checkError(err)

	fmt.Println(string(res))
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
