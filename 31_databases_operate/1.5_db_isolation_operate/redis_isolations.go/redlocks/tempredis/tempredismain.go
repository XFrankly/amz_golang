package tempredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/// 启动和停止redis server 临时进程

type Config map[string]string

func (c Config) Socket() string {
	return c["unixsocket"]
}

const (
	/// 当启动成功时，ready字符串表示 redis-server 字符串 打印标准输出，
	ready = "The server is now ready to accept conn"
)

// 封装配置，通过unixsocket启动 停止 单个redis-server 进程
type Server struct {
	dir       string
	config    Config
	cmd       *exec.Cmd
	stdout    io.Reader
	stdoutBuf bytes.Buffer
	stderr    io.Reader
}

// 启动随提供的配置信息 初始化一个新的 redis-server 进程
// redis-server 将监听一个本地临时unix socket 进程
/// 如果有任何原因 redis-server没有成功启动，一个错误将被返回
func Start(config Config) (server *Server, err error) {
	if config == nil {
		config = Config{}
	}

	dir, err := ioutil.TempDir(os.TempDir(), "tempredis")
	if err != nil {
		return nil, err
	}

	if _, ok := config["unixsocket"]; !ok {
		config["unixsocket"] = fmt.Sprintf("%s/%s", dir, "redis.sock")
	}
	if _, ok := config["port"]; !ok {
		config["port"] = "0"
	}

	server = &Server{
		dir:    dir,
		config: config,
	}
	fmt.Printf("server config:%+v\n", server)
	err = server.start()
	if err != nil {
		return server, err
	}

	// 阻塞直到redis 准备好连接
	err = server.waitFor(ready)

	return server, err
}

func (s *Server) start() (err error) {
	if s.cmd != nil {
		return fmt.Errorf("redis-server has already been started.")
	}
	s.cmd = exec.Command("redis-server", "-")

	stdin, _ := s.cmd.StdinPipe()
	s.stdout, _ = s.cmd.StdoutPipe()

	err = s.cmd.Start()
	if err == nil {
		err = writeConfig(s.config, stdin)
	}
	return err
}

func writeConfig(config Config, w io.WriteCloser) (err error) {
	for key, value := range config {
		if value == "" {
			value = "\"\""
		}
		_, err = fmt.Fprintf(w, "%s %s\n", key, value)
		if err != nil {
			return err
		}
	}
	return w.Close()
}

// 等待阻塞，直到redis-server 打印提供的string 到 stdout
func (s *Server) waitFor(search string) (err error) {
	var line string

	scanner := bufio.NewScanner(s.stdout)
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Fprintf(&s.stdoutBuf, "%s\n", line)
		if strings.Contains(line, search) {
			return nil
		}
	}

	err = scanner.Err()
	if err == nil {
		err = io.EOF
	}
	return err
}

/// 套接字 返回 完整路径到 本地redis-server 套接字
func (s *Server) Socket() string {
	return s.config.Socket()
}

/// 输出控制台 阻塞直到redis-server 返回 并且返回完整 输出到控制台输出
func (s *Server) Stdout() string {
	io.Copy(&s.stdoutBuf, s.stdout)
	return s.stdoutBuf.String()
}

/// 标准错误 阻塞，直到redis-server 返回 ，然后返回完整的输出到控制台
func (s *Server) Stderr() string {
	bytes, _ := ioutil.ReadAll(s.stderr)
	return string(bytes)
}

/// 优雅地关闭 redis-server， 它将返回一个错误，如果redis-server 终止失败
func (s *Server) Term() (err error) {
	return s.signalAndCleanup(syscall.SIGTERM)
}

/// kill 强制结束 redis-server 进程，它将返回一个错误，如果redis-server kill失败
func (s *Server) Kill() (err error) {
	return s.signalAndCleanup(syscall.SIGKILL)
}

func (s *Server) signalAndCleanup(sig syscall.Signal) error {
	s.cmd.Process.Signal(sig)
	_, err := s.cmd.Process.Wait()
	os.RemoveAll(s.dir)
	return err
}

// func main() {
// 	///// 使用redis 临时服务进程
// 	server, err := Start(Config{"databases": "9"})
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer server.Term()

// 	conn, err := redis.Dial("unix", server.Socker())
// 	defer conn.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	setrst, seterr := conn.Do("SET", "foo", "bar")
// 	fmt.Println("redis set foo bar:", setrst, seterr)
// 	getrst, geterr := conn.Do("GET", "foo")
// 	fmt.Println("redis set foo bar:", getrst, geterr)

// }
