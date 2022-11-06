package main

//HandleFunc 在 DefaultServeMux 中为给定模式注册处理函数。ServeMux 的文档解释了模式是如何匹配的。
import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	port := ":8085"
	h1 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #1!\n")
	}
	h2 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #2!\n")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/endpoint", h2)
	fmt.Printf("server start:%v\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
