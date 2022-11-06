package main

import (
	"fmt"
	"log"
	"net/http"
)

// 本文件服务显示全部文件 和 目录 包括 . 隐藏文件
func main() {
	// Simple static webserver:
	port := ":8090"
	fmt.Printf("file server start:%v\n", port)
	log.Fatal(http.ListenAndServe(port, http.FileServer(http.Dir("./"))))
}
