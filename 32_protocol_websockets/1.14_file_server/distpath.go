package main

//先绑定url path
// 在備用 URL 下提供磁盤 (/tmp) 上的目錄
// 路徑（/tmpfiles/），使用 StripPrefix 修改請求
// 在 FileServer 看到之前 URL 的路徑：
import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//

	port := ":8092"
	fmt.Printf("file server start:%v\n", port)
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("./tmp"))))
	log.Fatal(http.ListenAndServe(port, nil)) //http.FileServer(http.Dir("./tmp"))
	//http.FileServer(http.Dir("./"))))
}
