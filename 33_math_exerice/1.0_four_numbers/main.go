package main

import "github.com/aerospike/aerospike-client-go/types/rand"

/*
  任意 4个数字 组成的 整数（必须有一个大于6）
  它们的最大组合 减去 最小的组合 多次操作后，将收敛到 6437
*/

func gets_numbers() (int, int) {
	// 获取两个数字
	// a := rand.Int64()
	outcome := rand.Intn(6)
	return outcome
}
