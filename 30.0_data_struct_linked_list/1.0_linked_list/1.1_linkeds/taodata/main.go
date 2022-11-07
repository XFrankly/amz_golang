package taodata

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	DT = MakeNewDataTao("B", 8, 49)
)

func DoGua(env string, save bool) int {
	/*
	   执行一次运势 算卦
	   :return:
	*/
	fmt.Println(start)
	if env == "" {
		fmt.Println("占卜运数 环境,生月: 选择A,B,C对应哪月 A (一,二,三,四), B (五,六,七,八), C (九,十,冬,腊): ")
		fmt.Scan(&env)

	}
	newGua := MakeNewDataTao(env, 8, 49)
	ng, n := newGua.KanGuaOrigin()
	result, theNob := newGua.KanGua(ng, n)

	distriNumbs := DistributeDatas["T10000nonDiss"]
	ind := strconv.Itoa(theNob)
	distriNumb := distriNumbs[ind]

	fmt.Println("已测算此卦: 每一万人 与你相同的有 ", distriNumb, " 人")

	if save {
		fmt.Println("回车保存结果并退出:")
		var act = ""
		fmt.Scanln(&act)
		fmt.Println(".", act)
		fileName := "guali_" + time.Now().String()[:18]
		fileNames := strings.Replace(fileName, " ", "", -1)
		fileNames = strings.Replace(fileNames, ":", ".", -1)
		fileNames = strings.Replace(fileNames, "-", ".", -1) + ".txt"
		f, err := os.Create(fileNames)
		if err != nil {
			msg := fmt.Sprintf("Can not write data:%v\n", err)
			panic(msg)
		}
		f.WriteString(fileNames + "\n联系作者:hahamx@foxmail.com" + "\n")
		f.WriteString(result)
		f.Close()
		fmt.Println("已保存:", fileNames)
	}

	return theNob

}
