/* Go quine 自举 将把本模块代码内容打印*/
package main

import "fmt"

func main() {
	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)
}

var q = `/* Go quine */
package main

import "fmt"

func main() {
	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)
}
var q= `
