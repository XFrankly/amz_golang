package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

/*
http 完整示例 https://pkg.go.dev/net/http#pkg-examples
FileServer
FileServer (DotFileHiding)
FileServer (StripPrefix)
Get
Handle
HandleFunc
Hijacker
ListenAndServe
ListenAndServeTLS
NotFoundHandler
ResponseWriter (Trailers)
ServeMux.Handle
Server.Shutdown
StripPrefix
*/

func HttpFileServer() {
	//func FileServer(root FileSystem )
	//FileServer 返回一个处理 HTTP 请求的处理程序，文件系统的内容以根为根。
	//
	//作为一种特殊情况，返回的文件服务器会将任何以“/index.html”结尾的请求重定向到相同的路径，而不是最终的“index.html”。
	//
	//要使用操作系统的文件系统实现，请使用 http.Dir：
	// 启动一个 文件服务
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./"))))
}

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
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

func DotFileHidings() {
	fsys := dotFileHidingFileSystem{http.Dir(".")}
	http.Handle("/", http.FileServer(fsys))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HttpStripPrefixFileServer() {
	// To serve a directory on disk (/tmp) under an alternate URL
	// path (/tmpfiles/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func NewPeopleHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(w, "This is people handler.")
	})
}
func Http404Handler() {
	//	func NotFounHandler() Handler
	//	NotFoundHandler 返回一个简单的请求处理程序，该处理程序以“未找到 404 页面”回复每个请求。
	mux := http.NewServeMux()
	// 创建 简单 处理程序 返回 404
	mux.Handle("/res", http.NotFoundHandler())
	//	创建简单处理程序 返回 200
	mux.Handle("/res/people/", NewPeopleHandler())
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func HttpStripPrefix() {
	//	func StripPrefix(prefix string , h Handler ) Handler
	//StripPrefix 返回一个处理 HTTP 请求的处理程序，方法是从请求 URL 的路径（和 RawPath，如果设置）中删除给定的前缀并调用处理程序 h。
	//StripPrefix 通过回复 HTTP 404 not found 错误来处理对不以前缀开头的路径的请求。前缀必须完全匹配：
	//如果请求中的前缀包含转义字符，则回复也是 HTTP 404 not found 错误。
	mux := http.NewServeMux()
	mux.Handle("/filepath/", http.StripPrefix("/filepath/", http.FileServer(http.Dir("../apis"))))
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func HttpHijacker() {
	//type Hijacker interface {
	// // Hijack 让调用者接管连接。// 在调用劫持 HTTP 服务器库之后// 不会对连接做任何其他事情。// // 管理和关闭连接成为调用者的责任。//
	//返回的 net.Conn 可能已经设置了读取或写入期限 // 取决于服务器的配置// 服务器。调用者有责任根据需要设置// 或清除这些截止日期。
	//返回的 bufio.Reader 可能包含来自客户端的未处理的缓冲数据。// // 调用 Hijack 后，原来的 Request.Body 一定不能 使用。
	//原始请求的上下文保持有效并且// 在请求的 ServeHTTP 方法 返回之前不会被取消。
	//	Hijack() ( net . Conn , * bufio . ReadWriter ,错误)
	// http/1 支持， http/2不支持
	mux := http.NewServeMux()
	mux.HandleFunc("/hijack", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()
		bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
		bufrw.Flush()
		s, err := bufrw.ReadString('\n')
		if err != nil {
			log.Printf("error reading string: %v", err)
			return
		}
		fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
		bufrw.Flush()
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func HttpGetExample() {
	res, err := http.Get("http://www.google.com/robots.txt") //"http://127.0.0.1:8083/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("res, err, %s", body)
	fmt.Println(res, err)
}

func HttpResponseWriter() {
	//	HTTP 处理程序使用 ResponseWriter 接口来构造 HTTP 响应。
	//在 Handler.ServeHTTP 方法返回后，不能使用 ResponseWriter。
	//	// Header 返回将由WriteHeader 发送的头映射。Header map 也是// Handlers 可以设置 HTTP 尾部的机制。// // 在调用 WriteHeader（或// Write）之后更改标题映射无效，除非修改后的标题是// 拖车。// // 有两种设置 Trailers 的方法。首选方法是// 在标头中预先声明您稍后将发送的预告片 // 通过将“预告片”标头设置为稍后将出现的预告片键的名称。在这种情况下，这些// 标头映射的键被视为// 拖车。请参阅示例。第二种方式，拖车
	//	// 直到第一次写入之后，处理程序才知道键， // 是用 TrailerPrefix
	//	// 常量值作为 Header 映射键的前缀。
	//	请参阅预告片前缀。
	//	// 要禁止自动响应标头（例如“Date”），请将
	//	// 它们的值设置为 nil。
	//Header() Header
	//Write 将数据写入连接，作为 HTTP 回复的一部分。// // 如果还没有调用 WriteHeader，Write在写入数据之前调用 // WriteHeader(http.StatusOK)。如果 Header // 不包含 Content-Type 行，Write 会将 Content-Type 集// 添加到将写入数据的初始 512 字节传递给// DetectContentType 的结果。此外，如果所有写入的总大小

	// 数据小于几 KB 并且没有 Flush 调用，
	// Content-Length 标头会自动添加。
	//
	// 根据 HTTP 协议版本和客户端，调用
	// Write 或 WriteHeader 可能会阻止将来读取
	// Request.Body。对于 HTTP/1.x 请求，处理程序应 在写入响应之前读取任何 // 所需的请求正文数据。一旦
	// 标头被刷新（由于显式 Flusher.Flush
	// 调用或写入足够的数据以触发刷新），请求正文
	// 可能不可用。对于 HTTP/2 请求，Go HTTP 服务器允许
	// 处理程序继续读取请求正文，同时
	// 写入响应。但是，可能不支持此类行为
	// 所有 HTTP/2 客户端。如果可能，处理程序应在写入之前阅读
	// 以最大限度地提高兼容性。
	// Write([] byte) (int, error)
	//读取 Request.Body 时服务器自动发送的 100-continue 响应头除外。
	//WriteHeader(statusCode int)
	mux := http.NewServeMux()
	mux.HandleFunc("/sendstrailers", func(w http.ResponseWriter, req *http.Request) {
		// Before any call to WriteHeader or Write, declare
		// the trailers you will set during the HTTP
		// response. These three headers are actually sent in
		// the trailer.
		w.Header().Set("Trailer", "AtEnd1, AtEnd2")
		w.Header().Add("Trailer123", "AtEnd3")
		w.Header().Add("Authoration", "Bear amdkljdwihwdn1213kn12ned1lj2ni3n123jn123jn12kljn3")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		w.Header().Set("AtEnd1", "value 1")
		io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
		w.Header().Set("AtEnd2", "value 2")
		w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type apiHandler struct{}

// 重写ServeHTTP
func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func HttpServeMux() {
	//	func (mux * ServeMux ) Handle(模式字符串, handler Handler )
	//Handle 为给定模式注册处理程序。如果模式的处理程序已经存在，Handle 会发生恐慌。
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		w.Header().Set("Trailer", "AtEnd1, AtEnd2")
		w.Header().Add("Trailer123", "AtEnd3")
		w.Header().Add("Authoration", "Bear amdkljdwihwdn1213kn12ned1lj2ni3n123jn123jn12kljn3")
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func HttpGoodWayCloseServer() {
	//func (srv * Server ) 关机(ctx context . Context ) error
	//Shutdown 优雅地关闭服务器而不中断任何活动连接。关闭首先关闭所有打开的侦听器，然后关闭所有空闲连接，然后无限期地等待连接返回空闲状态，然后关闭。
	//如果提供的上下文在关闭完成之前过期，则 Shutdown 返回上下文的错误，否则返回关闭服务器的底层侦听器返回的任何错误。
	//当调用 Shutdown 时，Serve、ListenAndServe 和 ListenAndServeTLS 立即返回 ErrServerClosed。确保程序不会退出，而是等待 Shutdown 返回。
	//Shutdown 不会尝试关闭或等待被劫持的连接，例如 WebSockets。Shutdown 的调用者应单独通知此类长期连接的关闭并等待它们关闭（如果需要）。
	//有关注册关闭通知功能的方法，请参阅 RegisterOnShutdown。
	//一旦在服务器上调用了 Shutdown，它就不能被重用；将来对 Serve 等方法的调用将返回 ErrServerClosed。
	var srv http.Server
	// ctrl + c 优雅地开启 关闭服务
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		fmt.Println("server shutdown...")
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	fmt.Println("server start...")
	<-idleConnsClosed
}
func main() {
	//HttpFileServer()
	//DotFileHidings()
	//HttpStripPrefixFileServer()
	//Http404Handler()

	//HttpStripPrefix()
	//HttpHijacker()
	//HttpGetExample()

	//HttpResponseWriter()

	//HttpServeMux()
	HttpGoodWayCloseServer()
}
