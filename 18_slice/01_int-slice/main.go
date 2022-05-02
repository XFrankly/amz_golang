package main

import "fmt"

func main() {

	mySlice := []int{1, 3, 5, 7, 9, 11}
	fmt.Printf("%T\n", mySlice)
	fmt.Println(mySlice)

	mySlice1 := []int{0}
	mySlice2 := append(mySlice1, mySlice...) //arr_rpush(0, mySlice...) //append(mySlice1, mySlice...)
	fmt.Println(mySlice2)
	mySlice3 := mySlice2[:]
	for i := 0; i < len(mySlice2); i++ {
		if mySlice2[i] == 3 {
			mys1 := mySlice3[i+1:]
			mys2 := mySlice3[:i]
			// fmt.Println(mys1, mys2)

			mySlice4 := append(mys2, mys1...)
			fmt.Println(mySlice4)
		} else {
			continue
		}
	}

	fmt.Println(mySlice2, mySlice3)

	mySlice5 := []int{}
	for i := 0; i < len(mySlice2); i++ {
		//var v int
		//v := mySlice2[i] //mySlice2[i]
		mys5 := []int{mySlice2[i]}
		mySlice6 := append(mys5, mySlice5...)
		mySlice5 := append(mySlice6)
		fmt.Println(mySlice5)
	}
	fmt.Println(mySlice5)
}
