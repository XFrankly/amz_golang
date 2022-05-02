package main

//// 基于 空接口和反射的 容器，实现泛型

/// interface 这个泛型太 宽泛，使用反射进行类型检查

import (
	"fmt"
	"reflect"
)

type Container struct {
	// 通过传入存储元素类型 和 容量 来初始化 容器
	s reflect.Value
}

func NewContainer(t reflect.Type, size int) *Container {
	// 基于切片类型实现的容器，这里通过反射动态初始化这个底层切片
	return &Container{s: reflect.MakeSlice(reflect.SliceOf(t), 0, size)}
}

func (c *Container) Put(val interface{}) error {
	// 通过反射对 实际传递来的 元素类型进行运行时检查
	// 如果与容器初始化设置的元素类型不同，则返回错误信息
	// c.s.Type() 对应的是 切片类型，c.s.Type().Elem()对应的才是切片元素类型
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		// c.s 切片元素类型 与 传入参数不同
		return fmt.Errorf("put error:cannot put a %T into a slice of %s", c.s.Type().Elem())
	}
	// 如果类型检查通过则将其添加到容器
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
	return nil
}

func (c *Container) Get(val interface{}) error {
	// 还是通过反射对元素 类型进行检查，如果不通过则返回错误信息
	// kind 与 Type 相比范围更大，表示类别，如指针，而Type则对应具体类型，如 *int
	// 由于 val是指针类型，所有需要通过reflect.ValueOf(val).Elem() 获取指针指向的类型
	if reflect.ValueOf(val).Kind() != reflect.Ptr || reflect.ValueOf(val).Elem().Type() != c.s.Type().Elem() {
		return fmt.Errorf("get error:needs *%s but got %T", c.s.Type().Elem(), val)
	}
	// 将容器第一个索引位置值赋值给 val 指针
	reflect.ValueOf(val).Elem().Set(c.s.Index(0))
	// 然后删除容器第一个索引位置值
	c.s = c.s.Slice(1, c.s.Len())

	return nil
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	// 初始化容器，元素类型和nums中的元素类型相同
	c := NewContainer(reflect.TypeOf(nums[0]), 16)
	for _, n := range nums {
		if err := c.Put(n); err != nil {
			panic(err)
		}
		// 从容器读取元素，将返回结果初始化为0
		num := 0
		if err := c.Get(&num); err != nil {
			panic(err)
		}
		// 打印返回结果值
		fmt.Printf("%v, (%T)\n", num, num)
	}

	err := c.Put("s")     //put error:cannot put a *reflect.rtype into a slice of %!s(MISSING)
	err2 := c.Get("s100") //get error:needs *int but got string
	fmt.Println(err, err2)
}
