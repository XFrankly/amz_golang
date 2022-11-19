package main

import (
	"fmt"
)

func Exal(xs int) []error {
	aggreErr := AggreErrClosure()

	if xs > 5 {
		errs := fmt.Errorf("%v is over 5\n", xs)
		aggreErr(&errs)
	}

	if xs > 10 {
		errs := fmt.Errorf("%v is over 10\n", xs)
		aggreErr(&errs)
	}
	return aggreErr(nil)

}

// 闭包的示例
func AggreErrClosure() func(err *error) []error {
	var errs []error

	return func(err *error) []error {
		if err != nil {
			errs = append(errs, *err)
		}
		return errs
	}
}

func main() {
	err1 := Exal(1)
	err6 := Exal(6)
	err11 := Exal(11)
	fmt.Printf("err1:%v,\n err6:%v,\n errr11:%v\n", err1, err6, err11)
}
