package main

import (
	"fmt"
	"time"
)

func main() {
	var mp = map[string]int{
		"A": -1,
		"B": 0,
		"C": 1,
	}
	fmt.Println(mp["A"], mp["D"])

	msg := fmt.Sprintln("total", 12, "left", 9, "right", 3)
	fmt.Println(msg)

	//东非时区，与莫斯科时间一致
	Loc, _ := time.LoadLocation("Europe/Moscow") //"Asia/Shanghai")
	fp := time.Now().Local().In(Loc).String()
	fmt.Printf("date:%v, year:%v Month:%v  Day:%v\n", fp,
		time.Now().Local().In(Loc).Year(),
		time.Now().Local().In(Loc).Month(),
		time.Now().Local().In(Loc).Day())

	y, month, d := time.Now().Date()
	fmt.Println(y, month, d)
}
