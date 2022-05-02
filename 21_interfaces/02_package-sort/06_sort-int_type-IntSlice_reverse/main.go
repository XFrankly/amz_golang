package main

import (
	"fmt"
	"sort"
)

func main() {
	n := []int{7, 4, 8, 2, 9, 19, 12, 32, 3}

	fmt.Println(n)
	sort.Sort(sort.Reverse(sort.IntSlice(n)))  //降序
	fmt.Println(n)
	sort.Ints(n)   //升序
	fmt.Println(n)

}
