package main

import "fmt"

func Map[T any, V comparable](src []T, key func(T) V) map[V]T {
	var result = make(map[V]T)
	for _, v := range src {
		result[key(v)] = v
	}
	return result
}

/*
Map[T any, V comparable](src []T, key func(T) V) map[V]T {
    var result = make(map[V]T)
    for idx := range src {
        el := src[idx]
        result[key(el)] = el
    }
    return result
}
*/
var (
	SupField = []string{"name", "spec", "command"}
)

func ZipSlice[T any, V comparable](src []T, key func(T) V) map[V]T {
	var result = make(map[V]T)
	for idx := range src {
		el := src[idx]
		result[key(el)] = el
	}
	return result
}

func ZipSlices[T any, V comparable](src []T, key func(idx int) V) []map[V]T {
	var result = make(map[V]T)
	var resAll = []map[V]T{}
	total := len(src)
	size := len(SupField)
	if total < size {
		return resAll
	}
	for idx := 0; idx < total; idx += size {
		var j int
		for {
			if j >= size {
				j = 0
				break
			}
			result[key(j)] = src[idx+j]
			resAll = append(resAll, result)
			j++
		}
	}
	return resAll
}

func main() {

	slic2 := []interface{}{"python_cmd", "* * * 10 * *", "python -c \"import os; print(os.path())\"", "ls_daily", " * * * 10 * *", "ls -lah"}

	newMp := ZipSlices(slic2, func(idx int) string { return SupField[idx] })
	fmt.Printf("newMp:%v\n", newMp)
}
