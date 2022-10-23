package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

var (
	Logs = log.New(os.Stderr, "INFO -", 18)
)

type People struct {
	Age  int
	Name string
}

type SortByAge []People

func (a SortByAge) Len() int      { return len(a) }
func (a SortByAge) Swap(i, j int) { a[i].Age, a[j].Age = a[j].Age, a[i].Age }
func (a SortByAge) Less(i, j int) bool {
	return a[j].Age < a[j].Age
}
func (a SortByAge) DownSort() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Age > a[j].Age
	})
}
func (a SortByAge) UpSort() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Age < a[j].Age
	})
}
func (a SortByAge) NameUpSort() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Name < a[j].Name
	})
}

func NewPeoples(nums int) SortByAge {
	var ps []People
	ps = append(ps, People{Age: 25, Name: "Jack"})
	ps = append(ps, People{Age: 28, Name: "Lucy"})
	ps = append(ps, People{Age: 35, Name: "Frank"})
	ps = append(ps, People{Age: 36, Name: "Dave"})
	if nums > 4 {
		for i := 0; i < nums-4; i++ {
			nm := fmt.Sprintf("Jack%v", i)
			ps = append(ps, People{Age: 22 + i, Name: nm})
		}
	}
	return ps
}

func main() {
	ps := NewPeoples(10)
	Logs.Println(ps)
	Logs.Println(ps.Len())
	ps.Swap(0, 3)
	Logs.Println(ps)
	Logs.Println(ps.Less(0, 3))
	sort.Sort(SortByAge(ps))
	Logs.Println(ps)

	ps.DownSort()
	Logs.Println(ps)
	ps.UpSort()
	Logs.Println(ps)
	ps.NameUpSort()
	Logs.Println(ps)

	newSlice := []string{"a", "b", "c", "a1", "b1", "c1", "d1"}
	fmt.Println(newSlice[:3]) //[a b c]
	fmt.Println(newSlice[3:]) //[a1 b1 c1 d1]
	fmt.Println(newSlice[1:]) //[b c a1 b1 c1 d1]
	fmt.Println(newSlice[:1]) //[a]
}
