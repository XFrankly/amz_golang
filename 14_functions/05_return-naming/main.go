package main

import "fmt"

func main() {
	fmt.Println(greet("Jane ", "Doe"))
}

func greet(fname string, lname string) (s string) {  //在函数参数中定义 将返回s参数
	s = fmt.Sprint(fname, lname)
	return
}

/*
IMPORTANT
Avoid using named returns.


Occasionally named returns are useful. Read this article for more information:
https://www.goinggo.net/2013/10/functions-and-naked-returns-in-go.html
*/
