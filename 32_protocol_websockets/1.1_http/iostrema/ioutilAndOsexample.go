package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
os 的 示例不完整 完整请看这个 https://pkg.go.dev/os#pkg-examples
Chmod
Chtimes
CreateTemp
CreateTemp (Suffix)
Expand
ExpandEnv
FileMode
Getenv
LookupEnv
MkdirTemp
MkdirTemp (Suffix)
OpenFile
OpenFile (Append)
ReadDir
ReadFile
Unsetenv
WriteFile
*/
var (
	logg1   = log.New(os.Stderr, "[INFO] - ", 13)
	logger1 = log.New(os.Stderr, "[WARNING] - ", 13)
)

func ioutilReadAll() {
	r := strings.NewReader("Go is Great pramgraming language with systems programming in mind.")

	if b, err := ioutil.ReadAll(r); err == nil {
		logg1.Println(fmt.Sprintf("%s, %v", b, b)) // %s 显式 字符串， %v 显示结构和 存储值
		logg1.Printf("%T, %s", b, b)               // %s 显式 字符串， %v 显示结构和 存储值
		logg1.Println("b, err", b, err)
	} else {
		logger1.Println(err)
	}
}
func ioutilReadDir() {
	//	func ReadDir(dirname string ) ([] fs . FileInfo , error )
	//	ReadDir读取由 dirname命名的目录，并返回目录内容的 fs.FileInfo列表，按文件名排序。如果读取目录时发生错误，ReadDir将不返回任何目录条目和错误
	//	os.ReadDir时一个更有效和正确的选择，它返回fs.DirEntry 而不是 fs.FileInfo并且在读取目录 出错后仍返回 部分结果 需要go >1.16
	if files, err := ioutil.ReadDir("."); err == nil {
		for f, file := range files {
			logg1.Println("No file.", f, file.Name, file.Name())
		}
	} else {
		logg1.Println(err)
	}

}
func osReadDir() {
	//	func ReadDir(dirname string ) ([] fs . FileInfo , error )
	//	ReadDir读取由 dirname命名的目录，并返回目录内容的 fs.FileInfo列表，按文件名排序。如果读取目录时发生错误，ReadDir将不返回任何目录条目和错误
	//	os.ReadDir时一个更有效和正确的选择，它返回fs.DirEntry 而不是 fs.FileInfo并且在读取目录 出错后仍返回 部分结果 需要go >1.16
	if files, err := os.ReadDir("."); err == nil {
		for f, file := range files {
			logg1.Println("osReadDir No. file.", f, file.Name, file.Name())
		}
	} else {
		logg1.Println(err)
	}

}

func ioutilReadFile() {
	// 两个函数读取文件内容 ioutil.ReadFile， os.ReadFile
	//	func ReadFile(filename string)([]byte, error)
	//	ReadFile 读取由 filename 命名的文件并返回内容。成功的调用返回 err == nil，而不是 err == EOF。
	//	因为 ReadFile 读取整个文件，所以它不会将 Read 中的 EOF 视为要报告的错误。
	//从 Go 1.16 开始，此函数仅调用 os.ReadFile
	if content, errr := ioutil.ReadFile("./news.txt"); errr == nil {
		logg1.Println("read file news.txt", content, errr, fmt.Sprintf("%s", content))
	} else {
		logg1.Println(errr)
	}

	if cont, err := os.ReadFile("./news.txt"); err == nil {
		logg1.Println("os readfile", cont, err, fmt.Sprintf("%s", cont))
	} else {
		logg1.Println(err)
	}
}

func FsfileSystem(path interface{}) {
	var root string
	// os.DirFS 遍历指定路径的所有文件 名称
	if path == "" {
		root = "../"
	} else {
		root = fmt.Sprintf("%s", path)
	}
	logg1.Printf("loop name path :%v", root)
	fileSystem := os.DirFS(root)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logg1.Println(err)
		}
		logg1.Println(path)
		return nil
	})
}

func ioutilTempDir() {
	//	func TempDir(dir, pattern string) (name string, err error)
	//	TempDir 在目录 dir 中创建一个新的临时目录。目录名称是通过采用模式并在末尾应用随机字符串来生成的。如果模式包含“*”，则随机字符串替换最后一个“*”。
	//	TempDir 返回新目录的名称。如果 dir 是空字符串，则 TempDir 使用临时文件的默认目录（请参阅 os.TempDir）。
	//	多个程序同时调用 TempDir 不会选择同一个目录。不再需要时删除目录是调用者的责任。
	//	1.17 之后建议使用 os.MkdirTemp
	cont := []byte("temporary file's content")
	dir, err := ioutil.TempDir("./", "exampletest") //  如果不指定 dir ./  golang 将在 用户默认home目录创建 C:\Users\kukeg\AppData\Local\Temp\osexampletest3702339317
	if err == nil {
		logg1.Println("dir, err", dir, err, fmt.Sprintf("%v, %s", dir, dir))
	} else {
		logg1.Println(err)
	}
	diros, erros := os.MkdirTemp("./", "osexampletest")
	logg1.Println("os mkdir err rst:", erros)
	defer os.RemoveAll(dir) // cleanup
	defer os.RemoveAll(diros)
	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, cont, 0666); err == nil {
		logg1.Println("success write perm 0666", tmpfn, cont, fmt.Sprintf("%s", cont))
	} else {
		logg1.Println(err)
	}

}

func ioutilTempFile() {
	//func TempFile(dir, pattern string) (f *os.File, err error)
	//TempFile 在目录 dir 中创建一个新的临时文件，打开文件进行读写，并返回生成的 *os.File。
	//文件名是通过采用模式并在末尾添加一个随机字符串来生成的。如果模式包含“*”，则随机字符串替换最后一个“*”。
	//如果 dir 是空字符串，则 TempFile 使用临时文件的默认目录（请参阅 os.TempDir）。
	//多个程序同时调用 TempFile 不会选择同一个文件。调用者可以使用 f.Name() 来查找文件的路径名。不再需要时删除文件是调用者的责任。
	// go > 1.17 建议使用 os.CreateTemp
	cont := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "tempexample") // tmpfile.Name()  是带绝对路径的文件名
	tmpfilesuffix, errsuffix := ioutil.TempFile("", "suffixexample.*.txt")
	if err == nil {
		logg1.Println("tempfile", tmpfile, tmpfile.Name(), err)
	} else {
		logg1.Println(err)
	}
	if errsuffix == nil {
		logg1.Println("tmpfilesuffix, errsuffix", tmpfilesuffix.Name(), errsuffix)
	} else {
		logg1.Println(errsuffix)
	}
	ostmpefile, erros := os.CreateTemp("./", "ostempfile")
	if erros == nil {
		logg1.Println("ostmpefile", ostmpefile, ostmpefile.Name(), erros)
	} else {
		logg1.Println(erros)
	}
	defer os.Remove(tmpfile.Name()) //clean up
	defer os.Remove(ostmpefile.Name())
	defer os.Remove(tmpfilesuffix.Name())
	if wr, err := tmpfile.Write(cont); err == nil {
		logg1.Println("file write", wr, err)
	} else {
		logg1.Println(err)
	}
	if wrsf, err1 := tmpfilesuffix.Write(cont); err1 == nil {
		logg1.Println("suffix file write", wrsf, err1)
	} else {
		logg1.Println(err1)
	}
	if wros, err2 := ostmpefile.Write(cont); err2 == nil {
		logg1.Println("os file write", wros, err2)
	} else {
		logg1.Println(err2)
	}

	if err := tmpfile.Close(); err == nil {
		logg1.Println("file closed", tmpfile.Name())
	} else {
		logg1.Println("err close file", err)
	}
	if sferr := tmpfilesuffix.Close(); sferr != nil {
		logg1.Fatal(sferr)
	}
	if oserr := ostmpefile.Close(); oserr != nil {
		logg1.Fatal(oserr)
	}
	logg1.Println("file closed", tmpfilesuffix.Name(), ostmpefile.Name())
}

func ioutilWrite() {
	//	func WriteFile(filename string, data []byte, perm fs.FileMode) error
	//	WriteFile 将数据写入由文件名命名的文件。
	//	如果文件不存在，WriteFile 使用权限 perm（在 umask 之前）创建它；否则 WriteFile 在写入之前将其截断，而不更改权限。
	//	go > 1.6 推荐使用 os.WriteFile
	msg := []byte("Hello, GOlang")
	defer os.Remove("hello")
	if err := ioutil.WriteFile("hello", msg, 0664); err == nil {
		logg1.Println("ioutil write", err, msg, fmt.Sprintf("%v", err))
	} else {
		logg1.Println("ioutil write result", err, fmt.Sprintf("%v", err))
	}
	if err := os.WriteFile("hello", msg, 0664); err == nil {
		logg1.Println("os write", err, msg, fmt.Sprintf("%v", err))
	} else {
		logg1.Println("os write result", err, fmt.Sprintf("%v", err))
	}

}

func BufioReadLine() {
	// 按行读取文件
	var scanner *bufio.Scanner
	var files *os.File
	files, err := os.Open("news.txt")
	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner = bufio.NewScanner(files)
	scanner.Split(bufio.ScanLines)
	var text []string
	logg1.Println("len before read text", len(text))
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	defer files.Close()

	for x, each_in := range text {
		logg1.Println("x, each in ", x, each_in, "len:", len(text))
	}
}
func main() {
	//ioutilReadAll()
	//ioutilReadDir()
	//osReadDir()
	//
	//ioutilReadFile()
	//FsfileSystem(".")
	//ioutilTempDir()

	//ioutilTempFile()
	//ioutilWrite()
	BufioReadLine()
}
