package main

import "fmt"

func main() {
	intro := "Four 世"
	fmt.Printf("%T\n", intro) //查看类型
	fmt.Println(intro)        // 查看值
	bs := []byte(intro)       //转换为 byte类型 []uint8
	fmt.Println(bs)
	fmt.Printf("bs type %T\n", bs)
	fmt.Println("*********")
	fmt.Printf("%d\n", bs)

	for _, v := range bs {
		fmt.Printf("%d\t\t %#x\t %b\n", v, v, v)
	}
	fmt.Println("*********")
	y := 9999999999999999

	fmt.Printf("%d\t\t %#x\t %b\n", y, y, y)
	fmt.Println(&y)
	fmt.Sprint(y)
	fmt.Println("*********")

	z := 'h'
	fmt.Printf("%T\n", z)
}
