package stringutil

func reverseTwo(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
// 这说明了未导出的功能
// 可以被同一包中的导出函数使用
// this demonstrates how an unexported function
// can be used by an exported function in the same package
