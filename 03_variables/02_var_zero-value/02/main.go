package main

import (
	"bytes"
	"fmt"
	"os"
)

func ChangeBytesAndString() {
	//对于从字符串转换为字节切片，string -> []byte：
	str := "this is test"
	b1 := []byte(str)
	fmt.Printf("type of b1:%T type of str:%T \n", b1, str)
	fmt.Println(str, b1)

	//要将数组转换为切片，请[20]byte -> []byte：
	arr := [3]string{"1", "a", "c3"}
	ss1 := make([]int, 10)
	strings := []string{"shark", "cuttlefish", "squid", "mantis shrimp"}
	slices := arr[:]
	fmt.Printf("type of arr:%T type of slices:%T ss1:%T  strings:%T\n", arr, slices, ss1, strings)
	fmt.Println(arr, slices)

	//要将字符串复制到数组中，请执行以下操作 string -> [20]byte：
	strJson := "{\"user\": \"postgre\",\"password\": \"postgre.2022\"}"
	var arr2 = []byte(strJson)
	var arr21 [20]byte
	copy(arr21[:], "{\"user\": \"postgre\",\"password\": \"postgre.2022\"}")
	fmt.Printf("type of arr2:%T, %s \n", arr2, arr2)
	fmt.Println(arr2, arr21)

	//与上面相同，但首先将字符串显式转换为切片：
	var arr3 [30]byte
	copy(arr3[:], []byte(str))
	fmt.Printf("type of arr3:%T", arr3)
	fmt.Println(arr3)
	//内置函数仅从切片copy复制到切片。
	//数组是“底层数据”，而切片是“底层数据的视口”。
	//Using[:]使数组符合切片的条件。
	//字符串不符合 复制到的切片，但可以复制（字符串是不可变的）切片。
	//如果字符串太长，copy将只复制适合的部分字符串（然后可能只复制部分多字节符文，这将破坏结果字符串的最后一个符文）。
}
func main() {
	var a1 byte = 97
	var a2 byte = 98
	var a3 byte = 99
	var b bytes.Buffer // A Buffer needs no initialization.]
	var b3 []byte
	bb := b.Bytes()
	b.Write([]byte("\nHello "))
	fmt.Printf("%T", bb)
	fmt.Fprintf(&b, "world!\n")
	b.WriteTo(os.Stdout)

	// buf
	buf := bytes.Buffer{}
	buf.Write([]byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'})
	os.Stdout.Write(buf.Bytes())
	b3 = append(b3, a1)
	b3 = append(b3, a2)
	b3 = append(b3, a3)

	fmt.Printf("%T", b3)
	fmt.Println(b3)

	ChangeBytesAndString()
}
