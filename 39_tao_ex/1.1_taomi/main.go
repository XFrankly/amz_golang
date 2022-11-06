package main

import (
	"fmt"
	"math/rand"
	"time"
)

func DoRand() {
	rand.Seed(time.Now().UnixNano())
	top := 49 - 12
	sn := rand.Intn(top) + 12 // A random int in [12, 37]

	fmt.Println("sn: ", sn)
}
func main() {
	i := 0
	rand.Seed(time.Now().UnixNano())
	for {

		DoRand()
		r := rand.Intn(2)
		fmt.Println("r:", r)

		if i > 10 {
			break
		}
		i++
	}

}
