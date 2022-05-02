package main

import "fmt"
//bool:                    %t
//int, int8 etc.:          %d
//uint, uint8 etc.:        %d, %#x if printed with %#v
//float32, complex64, etc: %g
//string:                  %s
//chan:                    %p
//pointer:                 %p
//%b  格式化空格 所有有效的格式化动词
//（％b％e％E％e％f％F％g％G％x％X和％v）相等并接受
//十进制和十六进制表示法（例如：“ 2.3e + 7”，“ 0x4.5p-8”）
//和下划线分隔的下划线（例如：“ 3.14159_26535_89793”）。
func main() {
	for i := 1000000; i < 1000100; i++ {
		fmt.Printf("%d \t %b \t %x \n", i, i, i)
	}
}
