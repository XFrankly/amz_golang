package main

import "fmt"

/*
通用集合数据结构
类型参数的一个可能用例时实现一个通用的，类型安全的集合数据结构。
该实现非常简单，并且适合泛型
*/

/// Set 实现了由哈希表支持的通用集合数据结构。
//它不是线程安全的。
//// type comparable interface{ comparable }  类型可比较接口
/*
// 可比较是由所有可比较类型实现的接口
// (布尔值、数字、字符串、指针、通道、可比较类型的数组，
// 其字段都是可比较类型的结构）。
// 可比较接口只能用作类型参数约束，
// 不是变量的类型
*/

var (
	g = 0
)

type Set[T comparable] struct {
	Values map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] {
	m := make(map[T]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}
	return &Set[T]{
		Values: m,
	}
}

func (s *Set[T]) Add(values ...T) {
	for _, v := range values {
		///  占位符 struct{}{}
		// if s1 != nil {
		// 	s.Values[v] = *s1
		// } else {
		// 	s.Values[v] = struct{}{}
		// }
		s.Values[v] = struct{}{}

	}
}

type Mystring struct {
	Id    int
	Value string
}

func MakeString(a string) *Mystring {
	if a == "" {
		a = "default"
	}
	g += 1
	return &Mystring{
		Id:    g, /// 相当于 链表的 ID
		Value: a,
	}
}
func (s *Set[T]) Remove(values ...T) {
	/// 删除
	for _, v := range values {
		delete(s.Values, v)
	}
}

func (s *Set[T]) Contains(values ...T) bool {
	for _, v := range values {
		_, ok := s.Values[v]
		if !ok {
			return false
		}
	}
	return true
}
func (s *Set[T]) ToValues() []T {
	return s.toSlice()
}
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T](s.ToValues()...)
	for _, v := range other.ToValues() {
		if !result.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	///
	if s.Size() < other.Size() {
		return intersect(s, other)
	}
	return intersect(other, s)
}

//
func intersect[T comparable](smaller, bigger *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k, _ := range smaller.Values {
		if bigger.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

func (s *Set[T]) Size() int {
	return len(s.Values)
}

func (s *Set[T]) Clear() {
	s.Values = map[T]struct{}{}
}

func (s *Set[T]) String() string {
	//// 返回 列表切片的字符串形式
	return fmt.Sprint(s.toSlice())
}
func (s *Set[T]) toSlice() []T {
	result := make([]T, 0, len(s.Values))
	for k := range s.Values {
		result = append(result, k)
	}
	return result
}

//// 使用
func Do_generic() {
	s1 := NewSet(4, 4, -8, 15)
	s2 := NewSet("foo", "foo", "bar", "naz")
	fmt.Println(s1.Size(), s2.Size()) // 3, 3

	s1.Add(-16)
	s2.Add("hoge")
	fmt.Println(s1.Size(), s2.Size())                  // 4,4
	fmt.Println(s1.Contains(-16), s2.Contains("hoge")) // true, true

	s1.Remove(15)
	s2.Remove("naz")
	s2.Remove("baz") // Wrong
	fmt.Printf("values:%+v\n", s2.Values)
	fmt.Println(s1.Size(), s2.Size(), len(s1.Values), len(s2.Values)) // 3,3

	s3 := NewSet("hoge", "dragon", "fly")
	fmt.Println(s2.Union(s3).Size()) // 5
	fmt.Println(s2.Intersect(s3))    // [hoge]

	s1.Clear()
	s2.Clear()
	fmt.Println(s1.Size(), s2.Size()) // 0
}

func main() {
	/// Set 集合将数据保存在 map[T]struct{}{}
	// T是 类型参数，方法只需要将键放入，检索到映射并相互比较
	// 因此，允许这些操作的内置可比较约束对于类型参数T是足够的
	// 大多数代码与非泛型集 实现相同。
	/// 即使使用类型参数 推断 的使用示例 也没有 显示泛型类型位于其下方

	Do_generic()
}
