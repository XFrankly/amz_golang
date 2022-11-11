package main

import (
	"bstd/bstdata"
	"fmt"
)

func main() {
	// bstdata.TenTimes()  //初始算子值
	// bstdata.DothTest(49)
	var env string
	fmt.Println("占卜运数 环境,生月: 选择A,B,C对应哪月 A (一,二,三,四), B (五,六,七,八), C (九,十,冬,腊): ")
	fmt.Scan(&env)
	btm := bstdata.BstNew
	btm.KanGuaOrigin(env)

	// btm.Display()
}
