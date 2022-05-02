package main

import "fmt"

const (
	a = iota // 0
	b        // 1
	c        // 2
)

const (
	d = iota // 0
	e        // 1
	f        // 2
)

type ByteSize float64

const (
	a_          = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func main() {
	fmt.Println(a, b, c, d, e, f)
	//fmt.Println(b)
	//fmt.Println(c)
	//fmt.Println(d)
	//fmt.Println(e)
	//fmt.Println(f)
	fmt.Println(a_, KB, MB, GB, TB, PB, EB, ZB, YB)
}
