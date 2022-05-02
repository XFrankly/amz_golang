package main

import "fmt"
import "reflect"
import "errors"

func main() {
	data := []float64{43, 56, 87, 12, 45, 57}
	n := average(data)
	fmt.Println(n)
}
func In(haystack interface{}, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, errors.New("ErrUnSupportHaystack")
}
func average(sf []float64) float64 {
	// float
	fmt.Println(reflect.TypeOf(sf))  //类型查看
	fi := -1.2112
	fmt.Println(reflect.TypeOf(fi))  //类型查看

	// int
	var ai8 int8 = 127
	fmt.Println(reflect.TypeOf(ai8), ai8)  //类型查看
	ai8 = ai8 + 1
	fmt.Println(reflect.TypeOf(ai8), ai8)  //类型查看

	//数组
	coral := []string{"blue coral", "staghorn coral", "pillar coral"}  // 数组
	fmt.Println(coral) //显示全部
	fmt.Println(coral[1])  //显示第一个
	// 自写的函数 判断数组中是否存在某元素
	fmt.Println(In(coral, "blue coral1"))  //显示第一个
	// map的方式判断数组中 是否存在某元素, 依赖Go 中的 map 数据类型，通过 hash map 直接检查 key 是否存在
	//m := map[string]string{coral}
	//_,ok :=map("blue coral", coral)
	//if ok == true {
	//	fmt.Println(reflect.TypeOf(_), ok)
	//} else {
	//	fmt.Println("blue coral not exist ",   ok)
	//}

	// 浮点均值计算
	total := 0.0
	for _, v := range sf {
		total += v
	}
	return total / float64(len(sf))
}
