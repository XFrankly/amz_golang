package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	Loc, _ = time.LoadLocation("Asia/Shanghai")

	loc2, _ = time.LoadLocation("Local")
)

func main() {

	y := "2020"
	m := "6"
	result := GetMonthStartAndEnd(y, m)
	fmt.Println(result)
	fmt.Println(GetMonthTimesInt(y, m))
	/// 这一天属于第几周
	BelongWeeks("2022-09-16")

}

//GetMonthStartAndEnd 获取月份的第一天和最后一天
func GetMonthStartAndEnd(myYear string, myMonth string) map[string]time.Time {
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, _ := strconv.Atoi(myYear)

	timeLayout := "2006-01-02 15:04:05"
	theTime, _ := time.ParseInLocation(timeLayout, myYear+"-"+myMonth+"-01 00:00:00", Loc)
	newMonth := theTime.Month()

	t1 := time.Date(yInt, newMonth, 1, 0, 0, 0, 0, Loc)      //.Format(timeLayout)      //"2006-01-02")
	t2 := time.Date(yInt, newMonth+1, 0, 23, 59, 59, 0, Loc) //.Format(timeLayout) // "2006-01-02")
	result := map[string]time.Time{"start": t1, "end": t2}
	return result
}

func GetMonthTimesInt(myYear string, myMonth string) map[string]int {
	// 数字月份必须前置补零
	strEnds := GetMonthStartAndEnd(myYear, myMonth)
	result := make(map[string]int, 2)
	result["start"] = int(strEnds["start"].UnixMilli())
	result["end"] = int(strEnds["end"].UnixMilli())
	return result
}

//GetYearStartAndEnd 获取年份的第一天和最后一天 时间
func GetYearStartAndEnd(myYear string) map[string]time.Time {
	// 数字月份必须前置补零
	yInt, _ := strconv.Atoi(myYear)

	timeLayout := "2006-01-02 15:04:05"
	theTime, _ := time.ParseInLocation(timeLayout, myYear+"-"+"01"+"-01 00:00:00", Loc)
	newMonth := theTime.Month()

	t1 := time.Date(yInt, newMonth, 1, 0, 0, 0, 0, Loc) //.Format(timeLayout)      //"2006-01-02")
	t2 := time.Date(yInt, 12, 31, 23, 59, 59, 0, Loc)   //.Format(timeLayout) // "2006-01-02")
	result := map[string]time.Time{"start": t1, "end": t2}
	return result
}

//GetYearTimesInt 获取年份的第一天和最后一天 的时间戳
func GetYearTimesInt(myYear string) map[string]int {
	// 数字月份必须前置补零
	strEnds := GetYearStartAndEnd(myYear)
	result := make(map[string]int, 2)
	result["start"] = int(strEnds["start"].UnixMilli())
	result["end"] = int(strEnds["end"].UnixMilli())
	return result
}

//时间戳转字符时间
func TimeTempToStringDate(t int) string {
	timeNow := time.Unix(int64(t), 0)
	return timeNow.Format("2006-01-02 15:04:05")
}

//WeekIntervalTime 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func WeekIntervalTime(week int) map[string]time.Time {
	var (
		now    = time.Now()
		result = make(map[string]time.Time)
	)
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime := thisWeek.AddDate(0, 0, offset+7*week) //.Format("2006-01-02") + " 00:00:00"
	endTime := thisWeek.AddDate(0, 0, offset+6+7*week) //.Format("2006-01-02") + " 23:59:59"
	result["start"] = startTime
	result["end"] = endTime
	return result
}

//WeekIntervalTime 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func WeekTimeInt(week int) map[string]int {
	var (
		result = make(map[string]int)
	)
	rst := WeekIntervalTime(week)
	result["start"] = int(rst["start"].UnixMilli())
	result["end"] = int(rst["end"].UnixMilli())
	return result
}

//字符转时间戳  sec 秒，ms 毫秒，ns 纳秒
//秒(s), 毫秒(ms), 微秒(µs), 納秒(ns)
func StrDateToTimesTemp(s string, sec string) int {
	var unix_time int
	stringTime := s

	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", stringTime, Loc)

	if err == nil {
		if sec == "sec" {
			unix_time = int(the_time.Unix()) //1504082441

		} else if sec == "us" {
			unix_time = int(the_time.UnixMicro()) //1504082441999881
		} else if sec == "ns" {
			unix_time = int(the_time.UnixNano()) //1504082441999881

		} else {
			unix_time = int(the_time.UnixMilli()) //1504082441999

		}
	}
	return unix_time
}

// 当天时间所属那一年那一月第几周  ISO8601标准
func BelongWeeks(s string) {
	var (
		dateStr string
	)
	if s == "" {
		dateStr = "2019-12-31"

	}
	dateStr = s
	date, _ := time.Parse("2006-01-02", dateStr)
	// 获取当前时间数据当前第几周
	_, week := date.ISOWeek()

	weekday := date.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	// 获取日期所在周的周四日期
	thursday := date.AddDate(0, 0, int(4-weekday))
	// 所属月
	month := int(thursday.Month())
	// 所属年
	year := thursday.Year()
	// 获取日期所在周的周四所在的月是第几周
	_, week1 := time.Date(thursday.Year(), thursday.Month(), 4, 0, 0, 0, 0, time.Local).ISOWeek()
	// 日期所在周数减去日期所在那个周的周四的月份的4号所在的周数再加一即为本月第几周
	week = week - week1 + 1
	fmt.Println(year, month, week)

}

//任意时间字符 转换为时间格式
func StrToTime(s string) time.Time {
	return time.Now()
}
