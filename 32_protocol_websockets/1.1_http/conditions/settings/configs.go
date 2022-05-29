package settings

import "fmt"

var (
	Host     = "127.0.0.1"
	Port     = 18083
	Protocol = "http"
	Tls      = false

	Access_key = "1234567"
	Salt       = "&"
)

var (
	Secret_key  = fmt.Sprintf("%s%s", Access_key, Salt)
	AccessToken = "C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ==" // Token 不应写死 应该根据api的返回值 获取，然后保持在内存中
	Name        = "postgre"
	Passwords   = "postgre.2022"
)
