package main



import(
	"fmt"  //fmt使用与C的printf和scanf类似的功能实现格式化的 I/O
	"math"
)

func _diff_printf_printin() {
	m, n, p := 15, 25, 40

	fmt.Println( //PrintLn以默认格式指定
		"(m + n = p) :", m, "+", n, "=", p,
	)

	fmt.Printf(   //Printf根据指定的格式说明符进行格式设置
		"(m + n = p) : %d + %d = %d\n", m, n, p,
	)
}

func main() {
	_diff_printf_printin() // (m + n = p) : 15 + 25 = 40
							// (m + n = p) : 15 + 25 = 40
	_aaa()
	_natural()
	_logs()
	_floats()
}

func _aaa() {
	fmt.Println(math.Pow(2, 6))
}

func _natural() {
	//自然对数
	fmt.Println(math.Log(26))
}

func _logs() {
	fmt.Println(math.Log2(26))
	fmt.Println(math.Log10(26))
}

func _floats() {
	fmt.Println(math.Float64frombits(26)) //1.3e-322
	fmt.Println(math.Float64bits(26)) //4628011567076605952
	fmt.Println(math.Abs(-111)) //111
	print(math.Abs(-12222))
}
