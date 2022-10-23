package main

import (
	"fmt"
	"sort"
)

//逆序
func recover() {
	s := []int{5, 2, 6, 3, 1, 4} // unsorted
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	fmt.Println(s)
}

//浮点逆序
func recoverFloat() {
	s := []float64{5.1, 2.12, 6.2, 3.33, 1.1, 41.2} // unsorted
	sort.Sort(sort.Reverse(sort.Float64Slice(s)))
	fmt.Println(s)
}

//字符排序
func sortString() {
	s := []string{"Go", "Bravo", "Gopher", "Alpha", "Grin", "Delta"}
	sort.Strings(s)
	fmt.Println(s)
}
func main() {
	people := []struct {
		Name string
		Age  int
	}{
		{"Gopher", 7},
		{"Alice", 55},
		{"Vera", 24},
		{"Bob", 75},
	}
	sort.Slice(people, func(i, j int) bool { return people[i].Name < people[j].Name })
	fmt.Println("By name:", people)

	sort.Slice(people, func(i, j int) bool { return people[i].Age < people[j].Age })
	fmt.Println("By age:", people)

	recover()
	recoverFloat()
	sortString()
}
