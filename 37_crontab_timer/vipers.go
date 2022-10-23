package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	ConfigEnv         = "GCONFIG"
	ConfigDefaultFile = "config.yaml"
	ConfigTestFile    = "config.test.yaml"
	ConfigDebugFile   = "config.debug.yaml"
	ConfigReleaseFile = "config.release.yaml"
)

var (
	GDB     *sql.DB
	GDBList map[string]*sql.DB
	GREDIS  *redis.Client
	GCONFIG config.Server
	// GVP     *viper.Viper

	GLogs                            = InitLog() //                   = log.New(os.Stderr, "DEBUG -", 13)
	GTimers              timer.Timer = timer.NewTimerTask()
	GConcurrency_Control             = &singleflight.Group{}

	BlackCache local_cache.Cache
	lock       sync.RWMutex

	//登录 埋点是否加密 启用
	AuthHeader  string = "Bearer "
	PostEncrypt bool   = true
	AccessKey   string = "1234567890"
	//RPC 服务
	RedisAddr     = flag.String("redisAddr", "192.168.30.131:6379", "redis address")
	BasePath      = flag.String("basePath", "/rpcx_demo", "prefix path")
	RpcAddrr      = flag.String("addrr", "localhost:8973", "server address")
	RTSPLocalHost = "192.168.30.131"
)

//go:generate go-bindata -o=staticFile.go -pkg=packfile -tags=packfile ../resource/... ../config.yaml

func writeFile(path string, data []byte) {
	// 如果文件夹不存在，预先创建文件夹
	if lastSeparator := strings.LastIndex(path, "/"); lastSeparator != -1 {
		dirPath := path[:lastSeparator]
		if _, err := os.Stat(dirPath); err != nil && os.IsNotExist(err) {
			os.MkdirAll(dirPath, os.ModePerm)
		}
	}

	// 已存在的文件，不应该覆盖重写，可能在前端更改了配置文件等
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err2 := ioutil.WriteFile(path, data, os.ModePerm); err2 != nil {
			fmt.Printf("Write file failed: %s\n", path)
		}
	} else {
		fmt.Printf("File exist, skip: %s\n", path)
	}
}

func init() {
	for key := range _bindata {
		filePath, _ := filepath.Abs(strings.TrimPrefix(key, "."))
		data, err := Asset(key)
		if err != nil {
			// Asset was not found.
			fmt.Printf("Fail to find: %s\n", filePath)
		} else {
			writeFile(filePath, data)
		}
	}
}

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() {
				case gin.DebugMode:
					config = ConfigDefaultFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigDefaultFile)
				case gin.ReleaseMode:
					config = ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigReleaseFile)
				case gin.TestMode:
					config = ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigTestFile)
				}
			} else { // internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", ConfigEnv, config)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&GCONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&GCONFIG); err != nil {
		fmt.Println(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	GCONFIG.AutoCode.Root, _ = filepath.Abs("..")
	BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(GCONFIG.JWT.ExpiresTime)),
	)
	return v
}
