package main

import "fmt"

type vehicles interface{}

type vehicle struct {
	Seats    int
	MaxSpeed int
	Color    string
}

type car struct {
	vehicle //嵌入式结构体，也就是继承自 vehicle
	Wheels  int
	Doors   int
}

type plane struct {
	vehicle //嵌入式结构体，也就是继承自 vehicle
	Jet     bool
}

type boat struct {
	vehicle //嵌入式结构体，也就是继承自 vehicle，
	Length  int
}

func main() {
	v := vehicle{Seats: 4, Color: "red"}
	prius := car{vehicle: v}
	tacoma := car{vehicle: v}
	bmw528 := car{vehicle: v}
	boeing747 := plane{vehicle: v}
	boeing757 := plane{vehicle: v}
	boeing767 := plane{vehicle: v}
	sanger := boat{vehicle: v}
	nautique := boat{vehicle: v}
	malibu := boat{vehicle: v}
	rides := []vehicles{prius, tacoma, bmw528, boeing747, boeing757, boeing767, sanger, nautique, malibu}

	for key, value := range rides {
		fmt.Printf("k:%v - %v\n", key, value)
	}
}
