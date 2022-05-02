package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var logger = log.New(os.Stderr, "[INFO] --", 13)

type MyChans struct {
	// *MyChan //
	Read <-chan map[int]int //interface{} // 只读通道  为 channel 通道创建一个 按索引查看的方法
	//all   chan map[int]interface{}   // 可读可写
	Input chan<- map[int]int //interface{} // 只写通道
	// maxsize int

}

// MyChan的构造函数
func MakeChans(maxsize int) *MyChans {
	// 只读 只写 分开
	var MyC = make(chan map[int]int, maxsize)
	ret_mychan := &MyChans{
		Read:  MyC,
		Input: MyC,
		// maxsize: maxsize,
	}
	return ret_mychan
}

func IsBelongTo(i int, is []int) bool {
	/*
		i 是否 在 []int存在
	*/
	for j := range is {
		if i == j {
			return true
		}
	}
	return false
}
func IsPrime(x int) bool {
	/*
		x 是否质数
		偶数中除了2 都不是 素数，奇数的因数也没有偶数
	*/
	if x == 2 || x == 3 {
		return true
	} else if x%2 == 0 {
		return false // 排除偶数
	}
	sqrtX := int(math.Sqrt(float64(x)))
	for i := 3; i <= sqrtX; i += 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}

// func IsPrimeStep6(n int) []int {
// 	/*
// 			任何一个整数，总可以表示为
// 		    6n, 6n+1, 6n+2, 6n+3, 6n+4, 6n+5
// 		    6 进制
// 		    形容 6n+1 和 6n+5的数如果不是质素，它们的因素也将含有形如 6n+1 或 6n+5的数，因此可以得到

// 		    返回该数字的全部因子
// 	*/
// 	var all_prime = make([]int, 0) //[]int 创建一个切片长度为10，容量100

// 	if n == 2.0 || n == 3.0 {
// 		all_prime = append(all_prime, n)
// 	}
// 	nit := int(math.Sqrt(float64(n))) // ** 0.5
// 	fmt.Print(nit)
// 	for i := 5; i < nit+1; i += 6 {
// 		if n%i == 0 || (n%(i+2)) == 0 {
// 			// i 不属于 质数列表，并且 i 是质数
// 			if !IsBelongTo(i, all_prime) && IsPrime(i) == true {
// 				all_prime = append(all_prime, i)
// 			}
// 		}
// 	}
// 	// all_prime = append(all_prime, nit)
// 	return all_prime
// }

func bigger(a int) bool {
	if float64(a) >= math.Pow(10, 20) {
		return true
	}
	return false
}
func AllPrimeChildStep2(n int) []int {
	/*返回 n的全部质数因子*/
	result := make([]int, 0)
	x := int(math.Sqrt(float64(n))) //n - 1
	bigger := bigger(n)
	fmt.Println("bigger then Pow(10, 10):", bigger)
	if bigger == true { // 大于 1000 亿亿
		// for i := 5; i < x+1; i += 6 {
		// 	if n%i == 0 || (n%(i+2)) == 0 {
		// 		// i 不属于 质数列表，并且 i 是质数
		// 		if IsBelongTo(i, result) == false && IsPrime(i) == true {
		// 			result = append(result, i)
		// 		}
		// 	}
		// }
		fmt.Println("bigger result:", result)
		return result
	} else {
		for x > 1 {
			fmt.Println(n, x)
			if n%x == 0 {
				if IsPrime(int(n/x)) == false { // 合数拆分 为质数
					ss := AllPrimeChildStep2(int(n / x)) //  递归调用 到质数为止
					for _, s := range ss {
						if IsBelongTo(s, result) == false && IsPrime(s) == true {
							result = append(result, s)
						}
					}
				} else {
					result = append(result, int(n/x))
				}
				n = x
				x -= 1
			} else {
				x -= 1
			}
		}
		result = append(result, n)
	}

	return result
}
func First_n_prime_all_mul(n int) []int {
	/*
		按一定初始化的质数列表找出更多的质数
		    init prime list closer all number and its prime number
	*/
	init_prime := []int{2, 3, 5, 7, 11, 13, 17}
	failed := 0
	for i := 0; i < n; i++ {
		if failed >= 5 {
			return init_prime
		}
		if len(init_prime) > n {
			break
		}
		if i >= len(init_prime) {
			break
		}
		a := 1
		for _, j := range init_prime {
			a *= j
		}
		fmt.Println("a+1:", a+1)
		prime_list := AllPrimeChildStep2(a + 1)
		// if prime_list
		fmt.Println(prime_list)
		for _, p := range prime_list {
			if IsBelongTo(p, init_prime) == false && IsPrime(p) == true {
				init_prime = append(init_prime, p)
			}
		}
	}
	fmt.Println("init_prime length:", len(init_prime))
	return init_prime
}
func main() {
	t0 := time.Now()
	// baseN := 640336305379
	sp := AllPrimeChildStep2(8466443) //[12011 937 56897]

	fmt.Printf("sp:%+v", sp)
	// for i := range sp {
	// 	fmt.Println(i)
	// }
	// fmt.Println(36 >= math.Pow(10, 30)) //100000000000000000000000)

	// iprime := First_n_prime_all_mul(15)
	// fmt.Println(iprime)

	t := 0
	for i := 3; i <= 140000000; i += 2 { // 14亿以内的 所有 质数 2小时无法完成
		/*
			50分钟 完成计算 1 亿 4千万
			new prime: 139999973
			new prime: 139999991
			false
			cost time 48m56.7139038s

		*/
		if IsPrime(i) == true {
			fmt.Println("new prime:", i)
			t += 1
		}
	}
	fmt.Println("the number of prime less than baseN:", t)
	fmt.Println(128 > 493.23)
	fmt.Println("cost time", time.Since(t0))
}
