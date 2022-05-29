package resp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	logg   = log.New(os.Stderr, "[INFO] - ", 13)
	logger = log.New(os.Stderr, "[WARNING] - ", 13)
)

// RespMap 构造器
func InitRespMap() *RespMap {
	newRespMap := new(RespMap)
	return newRespMap
}

func MyDecoderMutilDict(jsonResps string, responsemap *RespMap) *RespMap {
	// 解析 多层 dict信息 解析并绑定为 目标结构体,
	//{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	if jsonResps == "" {
		jsonResps = `
	{"message":{"token":"C1WXDarU7CAEOVLfb6BclQqTMCmuHVkU7KbFTQ== ","user":"postgre"},"status":200}
	`
	}

	logg.Printf("jsonResps:%T", jsonResps)
	decResp := json.NewDecoder(strings.NewReader(jsonResps))
	for {
		// 解析 map
		//var Resp RespMap
		//if err := decResp.Decode(&Resp); err == io.EOF {
		if err := decResp.Decode(responsemap); err == io.EOF {
			break
		} else if err != nil { // 解析错误，格式不匹配
			logger.Println(err)
		}
		logg.Printf("%v, status:%d \n", responsemap.Message, responsemap.Status)
		logg.Printf("token: %s, user:%s \n", responsemap.Message["token"], responsemap.Message["user"])
	}
	return responsemap
}

func ReStrToStruct(loginResp []byte, respMap *RespMap) *RespMap {
	// 将 登录返回信息 转换为 相应的结构体
	logg.Printf("ReStrToStruct : %T \n", loginResp)
	loginRespStr := fmt.Sprintf("%s", loginResp)
	objresp := MyDecoderMutilDict(loginRespStr, respMap)

	logg.Printf("ReStrToStruct return:%T \n", objresp)
	return objresp
}
