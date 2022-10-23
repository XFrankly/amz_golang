package main

import (
	"fmt"
	"time"
)

func TimeStr() {
	// 获取当前时间
	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05.8156471 +0800 CST"))

	// 获取当前时间戳  unix 时间戳
	fmt.Println(t.Unix())

	// 时间 to 时间戳  设置时区时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(loc)
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05.8156471 +0800 CST", "2022-01-02 15:04:05.8156471", loc)
	fmt.Println(tt.Unix())
	// 时间戳 to 时间
	tm := time.Unix(1531293019, 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	/// 获取当前年月 时分秒
	y := t.Year()
	m := t.Month()
	d := t.Day()
	h := t.Hour()
	minute := t.Minute()
	seconds := t.Second()
	fmt.Println(y, m, d, h, minute, seconds)

}

func main() {
	TimeStr()
}
