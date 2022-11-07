package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// 本文件服务 不对外显示 隐藏目录
/*
containsDotFile 報告 name 是否包含以句點開頭的路徑元素。

	假定名稱由正斜杠分隔，保證 通過 http.FileSystem 接口。
*/
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

/*
dotFileHidingFile
是 dotFileHidingFileSystem 中使用的 http.File。

	用来保护隐藏目录 如 .git
	用來包裝http.File的Readdir方法，這樣我們就可以 從輸出中刪除以句點開頭的文件和目錄。
*/
type dotFileHidingFile struct {
	http.File
}

/*
Readdir 是嵌入 File 的 Readdir 方法的包裝器

	過濾掉所有以句點開頭的文件。
*/
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

/*
	dotFileHidingFileSystem

是一個隱藏的 http.FileSystem 隱藏的“點文件”被提供。
*/
type dotFileHidingFileSystem struct {
	http.FileSystem
}

/*
Open
是嵌入式 FileSystem 的 Open 方法的包裝器

	當 name 有文件或目錄時，提供 403 權限錯誤
	其名稱以路徑中的句點開頭。
*/
func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}

func main() {
	ports := ":8091"
	fsys := dotFileHidingFileSystem{http.Dir(".")}
	http.Handle("/", http.FileServer(fsys))
	fmt.Printf("start file server:%v\n", ports)
	log.Fatal(http.ListenAndServe(ports, nil))
}
