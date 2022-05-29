package resp

import (
	"net/http"
)

type HttpResponse struct {
	http.Response
}

type TokenUser struct {
	token string
	user  string
}

type RespMapStruct struct {
	// 多重使用 结构映射
	// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	Status  int
	message *TokenUser
}
type RespMap struct {
	// 只包含 一层，其他的使用 map
	// {"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}   // Support
	Status  int
	Message map[string]string
}

// 响应模块
type MyResp struct {
	Content          bool
	Content_consumed bool
	Next             interface{}

	//响应的 HTTP 状态的整数代码，例如404 或 200。 Integer Code of responded HTTP Status, e.g. 404 or 200.
	Status_code int

	//不区分大小写的字典 Response Headers.
	//#: 例如 ``headers['content-encoding']`` 将返回``'Content-Encoding'`` 相应头.
	Headers interface{}

	//#：响应的类似文件的对象表示（用于高级用法）。 使用“raw”需要在请求中设置“stream=True”。
	//#: 此要求不适用于请求的内部使用。
	Raw interface{}

	//最终到达网址的位置 Response.
	Url interface{}

	//#: 访问 r.text 时要解码的编码。
	Encoding interface{}

	//#: 列表对象 :class:`Response <Response>` 对象来自 请求的历史。
	// 任何重定向响应都将结束 上面这儿。该列表按从最早到最近的请求排序。
	History []string

	//#: 回复 原因 的文字版 HTTP Status, e.g. "Not Found" or "OK".
	Reason string

	//#: 服务器发回的 CookieJar。
	Cookies http.CookieJar //src.Cookies

	//#: 发送请求之间经过的时间量 和响应的到达（作为时间增量）。
	// 这个属性专门测量发送之间的时间 请求的第一个字节并完成对标头的解析。
	// 它 因此不受使用响应内容或 ``stream`` 关键字参数的值。
	Elapsed interface{} // datetime.timedelta(0)  0:00:00

	//#: 这个:class:`PreparedRequest <PreparedRequest>` 相应对象
	Request interface{}
}

// 构造 response
func ClientResponseControl(AuthToken string, ContentType string) *HttpResponse {
	return &HttpResponse{}
}
