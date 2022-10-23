package main

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

var (
	DivisionPrecision = 3
)

func Do_decimal() *int {
	/*
		加 Add   减 Sub  乘 Mul 除 Div
		decimal 运算得到的结果可以转成 自己想要的数据类型
		decimal.DivisionPrecision = 2  //保留两位数字小数，更多的四舍五入

		具有n位精度，通常表示在某个范围内，[10^k, 10^(k+1)] 其中k是整数，所有n位数字可以唯一识别
	*/
	decimal.DivisionPrecision = DivisionPrecision
	/// ExpMaxIterations 指定使用 ExpHullAbrham 方法计算精确自然指数值所需的最大迭代次数。
	var ExpMaxIterations = 1
	//// 将小数 JSON 编组为数字而不是字符串，则 MarshalJSONWithoutQuotes 应设置为 true。
	//警告：这对于多位数的小数是危险的，因为许多 JSON 解组器（例如：Javascript）
	//JS会将 JSON 数字解组为 IEEE 754 双精度浮点数，这意味着您可能会默默地失去精度
	var MarshalJSONWithoutQuotes = false
	//零常数，使计算更快。零不应直接与 == 或 != 进行比较，请改用 decimal.Equal 或 decimal.Cmp
	var zero = decimal.New(0, 1) //
	dec667 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
	fmt.Println(ExpMaxIterations, MarshalJSONWithoutQuotes, zero, dec667)

	// ddu1 := decimal.Decimal{}
	//// 0 值 对比是否相同
	fmt.Println(decimal.Zero.Equal(zero), decimal.Zero.Cmp(zero)) /// 与0对比结果 true，cmp结果 0
	decimal.DivisionPrecision = DivisionPrecision + 1
	price1 := decimal.NewFromFloat(3.14159268).Mul(decimal.NewFromInt(int64(100))) //.IntPart()
	price, err := decimal.NewFromString("136.02")                                  //136.02
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewFromInt(3) //3

	fee, _ := decimal.NewFromString(".035")       // 0.035
	taxRate, _ := decimal.NewFromString(".08875") // 0.08875

	subtotal := price.Mul(quantity) /// 乘法 price * quantity = 408.06

	preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1))) // 422.3421

	total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1))) //459.824961375
	fmt.Println(price1, price, quantity, fee, taxRate, subtotal, preTax)

	fmt.Println("Subtotal:", subtotal)                      // Subtotal: 408.06
	fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
	fmt.Println("Taxes:", total.Sub(preTax))                //37.482861375           // Taxes: 37.482861375
	fmt.Println("Total:", total)                            // Total: 459.824961375
	fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875

	/*
		浮点数，golang中只有 10进制 的，没有8进制 16进制
		常量 math.MaxFloat32 单精度 表示 float32能取到的 最大，约 3.4e38
		常量 math.MaxFloat64 双精度 表示 float64能取到的最大，约为 1.8e308
		最小的为 float32 和 float64 能表示的最小值分别为
		1.4e-45
		4.9e-324

	*/
	//// golang小数后 7 位是准确的
	//// 运算
	var myfloat01 float32 = 100000182
	var myfloat02 float32 = 100000187
	fmt.Println("myfloat:", myfloat01)
	fmt.Println("myfloat 2:", decimal.NewFromFloat32(myfloat02), decimal.NewFromFloat32(myfloat01), decimal.NewFromFloat32(5.0))
	after1Add := decimal.NewFromFloat32(myfloat01).Add(decimal.NewFromFloat32(5.0))
	f2Decimal := decimal.NewFromFloat32(myfloat02)
	fmt.Println(after1Add, f2Decimal, f2Decimal.Equal(after1Add))

	//// NewFromFloatWithExponent 将 float64 转换为 Decimal，具有任意数量的小数位数
	new_tax := total.Sub(preTax)
	fmt.Println(new_tax, decimal.NewFromFloatWithExponent(37.482861375, -3)) /// 取3位小数
	v, _ := new_tax.Value()                                                  /// 取值
	ia, _ := strconv.Atoi(v.(string))                                        // string 转换为 int
	strIa := strconv.Itoa(ia)                                                // int 转为 string
	fmt.Printf("%T, %T\n", ia, strIa)
	fmt.Println("\n", ia, strIa)
	//// NewFromFloat 通常有 15位的 精度
	fmt.Println(decimal.NewFromFloat(123.123123123123).String())
	fmt.Println(decimal.NewFromFloat(.123123123123123).String())
	fmt.Println(decimal.NewFromFloat(-1e13).String())
	//// NewFromFloat32 靠往返的浮点数中表示的有效数字的数量。这通常是 6-8 位，具体取决于输入
	fmt.Println(decimal.NewFromFloat32(123.123123123123).String())
	fmt.Println(decimal.NewFromFloat32(.123123123123123).String())
	fmt.Println(decimal.NewFromFloat32(-1e13).String())

	return nil
}

func main() {
	Do_decimal()
}
