package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	msg := "\n Do not dwell in the past, do not dream of the future, concentrate the mind on the present.\n"
	msg2 := "\n不要沉迷于过去，不要梦想未来，将精力集中在当下。\n"
	rdr := strings.NewReader(msg)
	rdrzh := strings.NewReader(msg2)
	io.Copy(os.Stdout, rdr)
	io.Copy(os.Stdout, rdrzh)

	rdr2 := bytes.NewBuffer([]byte(msg))
	io.Copy(os.Stdout, rdr2)

	res, _ := http.Get("http://www.google.com")
	//// 将 res.Body 的内容 copy 到 os.Stdout，也就是打印到控制台
	io.Copy(os.Stdout, res.Body)
	res.Body.Close()
	fmt.Println("%+v\n", res)
}
