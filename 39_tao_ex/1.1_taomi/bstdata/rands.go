package bstdata

import (
	"fmt"
	"math/rand"
	"time"
)

// 默认从 4 到 size 区间
func DoRand(size int) int {
	rand.Seed(time.Now().UnixNano())
	top := size - 4
	sn := rand.Intn(top) + 4 // A random int in [4, 45]

	fmt.Println("new sn: ", sn)
	return sn
}

// 默认从 lowst 到 tops 区间
func DoRange(lowst, tops int, env string) int {
	if tops <= lowst {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	top := tops - lowst
	sn := rand.Intn(top) + lowst // A random int in [lowst, 45]

	fmt.Println("new sn: ", sn)
	er := EnvRange[env]
	if er != 0 {
		if er == -1 {
			sn -= 4
		} else if er == 1 {
			sn += 4
		}
	}

	if sn < lowst+4 {
		sn = lowst + 4
	}
	if sn > tops-4 {
		sn = tops - 4
	}
	return sn
}

func TenTimes() {
	i := 0
	rand.Seed(time.Now().UnixNano())
	for {

		DoRand(49)
		r := rand.Intn(2)
		fmt.Println("r:", r)

		if i > 10 {
			break
		}
		i++
	}
}
