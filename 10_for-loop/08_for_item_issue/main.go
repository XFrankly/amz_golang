package main

import "fmt"

type Item struct {
	Name string
}

func main() {
	var all = []*Item{}
	var Items = []Item{Item{Name: "Jack"}, Item{Name: "Lucy"}}
	for _, item := range Items {
		// item := item
		all = append(all, &item)
	}
	fmt.Printf("items:%#v\n", Items)

	for ind, it := range all {
		fmt.Printf("ind:%#v it:%#v\n", ind, it)
	}

}
