package main

import (
	"fmt"
)

type Number interface {
	int | float64
}

// SliceFn implements sort.Interface for a slice of T.
// SliceFn 为 T 的一个切片实现 sort.Interface。
type SliceFn[T any] struct {
	s    []T
	less func(T, T) bool
}

func (s SliceFn[T]) Len() int {
	return len(s.s)
}
func (s SliceFn[T]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}
func (s SliceFn[T]) Less(i, j int) bool {
	return s.less(s.s[i], s.s[j])
}

// // SortFn sorts s in place using a comparison function.
// func SortFn[T any](s []T, less func(T, T) bool) {
// 	sort.Sort(SliceFn[T]{s, cmp})
// }

type A[N Number] struct {
	v1 N
	v2 N
}

func MakeAs[n Number](v n, v2 n) *A[n] {
	return &A[n]{v1: v, v2: v2}
}
func (b *A[n]) Adds(v *A[n]) {
	// 使用指针传递 则不必每次返回 传入对象，而是直接修改指针指向的对象的属性
	b.v1 += v.v1
	b.v2 += v.v2
}

type B struct {
	v1 int
	v2 int
}

func (b *B) Add(v *B) {
	// 使用指针传递 则不必每次返回 传入对象，而是直接修改指针指向的对象的属性
	b.v1 += v.v1
	b.v2 += v.v2
}

func main() {
	b1 := &B{v1: 11, v2: 22}
	v1 := &B{v1: 12, v2: 23}
	b1.Add(v1)
	fmt.Printf("b1:%+v\n", b1)
	fmt.Printf("v1:%+v\n", v1)

	//#################
	a1 := MakeAs(2, 232)
	va := MakeAs(12, 201)
	a1.Adds(va)
	fmt.Println(a1)
	fmt.Printf("va:%+v\n", va)
}
