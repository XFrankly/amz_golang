package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SeekStart   = 0 // seek relative to the origin of the file
	SeekCurrent = 1 // seek relative to the current offset
	SeekEnd     = 2 // seek relative to the end
)

var (
	logg   = log.New(os.Stderr, "[INFO] - ", 13)
	logger = log.New(os.Stderr, "[WARNING] - ", 13)
	//EOF 是当没有更多可用输入时 Read 返回的错误。
	//（读取必须返回 EOF 本身，而不是包装 EOF 的错误，因为调用者将使用 == 测试 EOF。）函数应该返回 EOF 仅表示输入的优雅结束。
	//如果 EOF 在结构化数据流中意外发生，则相应的错误是 ErrUnexpectedEOF 或提供更多详细信息的其他错误。
	EOF = errors.New("EOF")

	//ErrClosedPipe 是用于对封闭管道进行读取或写入操作的错误。
	ErrClosePipe = errors.New("io:read/write on closed pipe")

	//当对 Read 的许多调用都未能返回任何数据或 错误时，Reader 的某些客户端会返回 ErrNoProgress，这通常是 Reader 实现损坏的标志。
	ErrNoProgress = errors.New("multiple Read calls return no data or error")

	// ErrShortBuffer 意味着读取需要比提供的更长的缓冲区。
	ErrShortBuffer = errors.New("short buffer")

	// ErrShortWrite 表示写入接收的字节 少于请求的字节，但未能返回显式错误
	ErrShortWrite = errors.New("short write")

	// ErrUnexpectedEOF 表示在固定大小的块 或数据结构过程中遇到了 EOF
	t0 = time.Now()
)

func io_copy() {
	//	/ 把副本 从 src 复制到 dst，直到src上达到EOF或发生错误，它返回复制的字节数和复制时遇到的第一个 错误(如果有)，
	// 成功Copy 返回 err == nil， 如果 err == EOF 表示失败
	//	因为Copy被定义为从sr读取直到EOF 所以它不会将Read的EOF视为要报告的错误
	//	如果src实现 WriteTo接口，则通过调用src.WriteTo(dst)实现复制，否则，如果dst实现了ReaderFrom接口，则通过调用dst.ReadFrom(src)实现复制\
	// r 类型 *strings.Reader
	logg.Printf("####################io_copy\n")
	r := strings.NewReader("some io.Reader stream to be read \n")
	//var newr string
	// func io.Copy(dst Writer, src Reader) (written int64, err error) {
	if cr, err := io.Copy(os.Stdout, r); err != nil {
		logger.Fatal(err)
	} else if newcr, err := io.Copy(log.Writer(), r); err == nil {
		logg.Printf("%T", r)

	} else {
		fmt.Println("copy r:", r, "\ncost time:", time.Since(t0))
		logger.Println("copy r:", r, "\ncr:", cr, "\nnewcr:", newcr, "\ncost time:", time.Since(t0))
	}
}

func io_copy_buffer() {
	//func CopyBuffer(dst Writer , src Reader , buf [] byte ) (Writer int64 , err error )
	//	CopyBuffer 与 Copy相同，只是它通过提供的缓冲区进行暂存，而不是分配临时缓冲区
	//	如果buf 为 nil，则分配一个，否则如果长度为0，CopyBuffer发生混乱
	//	如果src实现 WriterTo 或 dst 实现 ReaderFrom 则不会使用 buf 执行复制
	logg.Printf("####################io_copy_buffer\n")
	r1 := strings.NewReader("first reader \n")
	r2 := strings.NewReader("second reader \n")
	buf := make([]byte, 8) //  不能为0
	// buf is 在这里使用
	if rb, err := io.CopyBuffer(os.Stdout, r1, buf); err != nil {
		logger.Println("rb:", rb)
		logger.Fatal(err)
		//} else if lgb, err := io.CopyBuffer(log.Writer(), r1, buf); err == nil {
		//	logg.Println("lgb:", lgb, "r1:", r1, "buf:", buf, err)
	} else {
		logger.Println("rb:", rb, err)
	}

	// 重用 不需要 allocate 一个 新的 扩展buffer
	if r2b, err := io.CopyBuffer(os.Stdout, r2, buf); err == nil {
		logg.Println("copy buffer r2b:", r2b, "r2:", r2)
	} else {
		logger.Println(err)
	}
}

func io_copy_Nbytes() {
	// 从 src 复制 N 位比特 到 dst
	//	func CopyN(dst Writer , src Reader , n int64 ) (写入int64 , err error )
	//	它返回复制的字节数和复制时遇到的最早错误。返回时，写入 == n 当且仅当 err == nil。
	//如果 dst 实现了 ReaderFrom 接口，则使用它实现副本。
	logg.Printf("####################io_copy_Nbytes\n")
	r := strings.NewReader("some io.Reader stream to be read.")
	if rn, err := io.CopyN(os.Stdout, r, 4); err == nil {
		logg.Println("os stdout rn", rn, "err", err, "\nr", r)
	}
	lw := logg.Writer()
	if rnl, err := io.CopyN(lw, r, 4); err == nil {
		logg.Println("logg writer result rnl:", rnl, "err", err, *&lw)
	}
}

func io_writer_pipe() {
	//	管道 Pipe 创建一个同步的内存管道。 可用于连接期望io.Reader的代码和期望 io.Writer代码
	//	管道的读取 和 写入 是 1:1 匹配的，除非需要多个读取来消耗单个写入，也就是，对 PipeWriter 的每次写入都会阻塞
	//	直到它满足 来自 PipeReader 的一个或多个读取，这些读取完全消耗了写入的数据
	//	数据直接从Write 复制到对应的Read 或 Reads
	//	没有内部 缓冲，彼此并行 或 使用 Close调用 的Read 和 Write是安全的，对Read的并行调用和对 Write的并行调用也是安全的
	//	各个调用将 顺序进行 门控
	logg.Printf("####################io_writer_pipe\n")
	r, w := io.Pipe() // 创建一个 一 一 对应的管道，r 读取，w写入
	logg.Println("Before w pipe r:", r, "w:", w)
	go func() {
		logg.Println(w, "some io.Reader stream to be read \n")
		w.Close()
	}()
	if rpipe, err := io.Copy(os.Stdout, r); err == nil {
		logg.Println("rpipe", rpipe, err, "\nr, w and os.stdout:", &r, &w, os.Stdout)
		logg.Println("r w:", r, w)
	}
}

func io_read_all() {
	//	func ReadAll(r Reader ) ([] byte , error )
	//	ReadALl 从r读取，直到出现错误 或 EOF 并返回它读取的数据。 成功的调用返回 err == nil
	//	而不是err == EOF。 因为ReadALL 被定义为从src读取到 EOF， 所以它不会将 Read 的EOF视为要报告的错误
	logg.Printf("####################io_read_all\n")
	r := strings.NewReader("Go is a general-purpose language designed with systems programming in mind.")
	if b, err := io.ReadAll(r); err == nil {
		logg.Printf("r type: %T", r)
		logg.Println("read r:", r, err, "b:", b, len(b))
	}

}

func io_read_all_from_api_booking() {
	// 从自定义接口 读取并返回数据
	// Reader  作为数据源
	r := strings.NewReader("{\"id\": \"1\",\"title\": \"London\",\"artist\": \"PostBettyCarter\",\"price\": 79.99}") // post 的传入body的数据
	resp, err := http.Post("http://127.0.0.1:8083/v1/booking/add", "applicaation/json", r)
	// req, err := http.NewRequest("POST", "http://127.0.0.1:8083/v1/booking", nil)
	// 关闭 响应包Body defer resp.Body.Close()
	defer resp.Body.Close()
	body, err1 := io.ReadAll(resp.Body) // 必须io read
	logg.Printf("%T \n", body)
	body_str := fmt.Sprintf("%s", body)
	logg.Printf("POST resp type %T, resp.Body %T, %T, %s \n", resp, &resp.Body, body_str, body)
	logg.Println("POST /v1/booking resp: ", resp, "\nBody:", body_str, len(body), "\n", *resp, "\nerr:", err, err1)

	r2 := strings.NewReader("")
	respc, _ := http.Post("http://127.0.0.1:8083/v1/booking/clean", "applicaation/json", r2)
	defer respc.Body.Close()
	bodyc, errc1 := io.ReadAll(resp.Body) // 必须io read
	bodyc_str := fmt.Sprintf("%s", bodyc)

	logg.Printf("%T bodyc_str:%s, Done. \n", bodyc, bodyc_str)
	logg.Printf("POST bodyc type %T, resp.Body %s,%v \n", bodyc, bodyc_str, errc1)
}

func io_read_at_least() {
	//func ReadAtLeast(r Reader , buf [] byte , min int ) (n int , err error )
	//ReadAtLeast 从r读取到 buf 直到它至少读取了 min 字节，它返回复制的字节数，如果读取的字节数少，则返回错误
	//	仅当未读取任何字节时，该错 为 EOF
	//	在读取少于min字节后 发生 EOF， ReadAtLeast 返回 err UnexpectedEOF，如果 min 大于buf的长度，ReadAtLeast返回 ShortBuffer 错误
	//	返回时 n > = min 当且仅当 err == nil 时。 如果 r 返回读取至少min 字节错误，则删除该错误
	r := strings.NewReader("some io.Reader stream to be read \n")
	buf := make([]byte, 14)

	logg.Println("before buf r.Len(:", r.Len())
	if rst, err := io.ReadAtLeast(r, buf, 4); err == nil {
		logg.Printf("buf14: %v %s total remain:%v \n", rst, buf, r.Len())
	} else {
		//logg.Fatal(err)
		logg.Printf("err 14: %v, %s", rst, err)
	}
	// buff smaller than minimal read size buff 空间小于 读取的字节 4
	// 读取失败 不影响 阅读器中的内容
	shortBuf := make([]byte, 3)
	if rst2, err := io.ReadAtLeast(r, shortBuf, 4); err == nil {
		logg.Printf("shortBuf:%v %s \n", rst2, err)
	} else {
		//logg.Fatal(err)
		logg.Printf("err s: %v, %s", rst2, err) //  0, short buffer

	}
	// minimal read size bigger than io.Reader stream buf空间 等于 读取的字节数 但是大于 可 可读取字节
	logg.Println("before read:", r.Len())
	bigBuf := make([]byte, 64)
	if brst, err := io.ReadAtLeast(r, bigBuf, 19); err == nil { // 如果 min大于 r.len() 返回 r的长度 20, unexpected EOF
		logg.Printf("bigBuf: %v, %v, %v, remain:%v", brst, err, bigBuf, r.Len()) // 20, unexpected EOF
		logg.Println(r)
	} else {
		//logg.Fatal(err)
		logg.Printf("err B: %v, %s", brst, err) // err B: 20, unexpected EOF

	}
}

func io_read_at_least_util_empty() {
	//func ReadAtLeast(r Reader , buf [] byte , min int ) (n int , err error )
	//ReadAtLeast 从r读取到 buf 直到它至少读取了 min 字节，它返回复制的字节数，如果读取的字节数少，则返回错误
	//	仅当未读取任何字节时，该错 为 EOF
	//	在读取少于min字节后 发生 EOF， ReadAtLeast 返回 err UnexpectedEOF，如果 min 大于buf的长度，ReadAtLeast返回 ShortBuffer 错误
	//	返回时 n > = min 当且仅当 err == nil 时。 如果 r 返回读取至少min 字节错误，则删除该错误
	r := strings.NewReader("some io.Reader stream to be read\n")
	buf := make([]byte, 5)
	const mins = 4
	logg.Println("before buf r.Len(:", r.Len())
	for {
		if r.Len() > 0 {
			if r.Len() < mins {
				if rstA, errA := io.ReadAll(r); errA == nil { //
					logg.Printf("buf14: %v %v total remain:%v \n", rstA, errA, r.Len())
					return
				} else {
					logg.Fatal("ReadAll err when len", r.Len())
				}
			}
			if rst, err := io.ReadAtLeast(r, buf, mins); err == nil { // 至少读4个
				logg.Printf("buf14: %v %s total remain:%v \n", rst, buf, r.Len())
			} else {
				//logg.Fatal("err 14: %v, %s", rst, err)
				logg.Printf("err 14: %v, %s", rst, err) // err 14: 3, unexpected EOF   读完3个后 报错 EOF，同时 阅读器内 已经 取完了
			}
		} else {
			logg.Fatal("err fatal: length:", r.Len())
		}
	}
}

func io_read_full_buf() {
	//	func ReadFull(r Reader , buf [] byte ) (n int , err error )
	//	io.ReadFull 读满buf, 它发挥复制的字节数，如果读取的字节数较少，则返回错误，
	//需要读满，但是 仅当未读取任何字节时，该错误为EOF 将返回 UnexpectedEOF
	//	当 n == len(buf) 当且仅当 err == nil， 返回至少读取 len(buf) 字节错误，则删除该错误
	r := strings.NewReader("some io.Reader stream to be read\n")

	buf := make([]byte, 4)
	if fr, err := io.ReadFull(r, buf); err == nil {
		logg.Println("read full:", fr, err, buf, r.Len())
	} else {
		logg.Println(err)
	}
	//  buf size 大于 可读 字节数
	longBuf := make([]byte, 64)
	if lf, err := io.ReadFull(r, longBuf); err == nil {
		logg.Println("longBuf:", lf, err, longBuf, r.Len())
	} else {
		logg.Println("err longBuf read", err, r.Len()) // 还剩29个字符
	}
}

func io_write_string() {
	// 使用io.WriteString 在控制台输出 字符
	if iw, err := io.WriteString(os.Stdout, "\nHello, World\n"); err == nil {
		logg.Println("io write", iw, err)
	} else {
		logg.Fatal("io writer err", err)
	}
	if iw, err := io.WriteString(logg.Writer(), "\nHello, World logg\n"); err == nil {
		logg.Println("logg io write", iw, err)
	} else {
		logg.Fatal("logg io writer err", err)
	}

}

// io.Reader 结构体
//type Reader interface {
//	Read(buf []byte) (n int, err error)
//}

func main() {
	//io_copy()
	fmt.Printf("start time now:%v \n", time.Now())

	//io_copy_buffer()
	//
	//io_copy_Nbytes()
	//io_writer_pipe()

	//io_read_all()

	//io_read_at_least()

	io_read_at_least_util_empty()
	// io_read_full_buf()

	// io_write_string()
	io_read_all_from_api_booking()
	logg.Println("end copy io buffer cost time:", time.Since(t0))
}
