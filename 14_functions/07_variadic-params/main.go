package main

import "fmt"

func main() {
	n := average(43, 56, 87, 12, 45, 57)
	fmt.Println(n)

	s := listdoub(11, 22, 33, 44, 55)
	fmt.Println(s)
}

func average(sf ...float64) float64 { //sf 为列表，可变的参数
	fmt.Println(sf)
	fmt.Printf("%T \n", sf) //显示格式
	var total float64
	for _, v := range sf {
		total += v
		fmt.Println(total, v, "in", sf)
	}
	return total / float64(len(sf)) // 返回均值
}

func listdoub(sf ...int32) int32 { //sf 为列表，可变的参数
	fmt.Println([]int32(sf))
	fmt.Printf("%T \n", sf) //显示格式

	var total int32
	for _, v := range sf {
		total += v
		fmt.Println(total, v, "in", sf)
	}
	return total / float64(len(sf)) // 返回均值
}
