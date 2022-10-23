package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(log.Writer(), "INFO ", 13)

func reverseTwo(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func count(n float64) {

	var ms interface{}
	for j := 1.0; j <= n; j++ {
		// j 的逆序值
		sj := fmt.Sprintf("%s", reverseTwo(fmt.Sprintf("%v", j)))
		//  j 的平方根
		ms = math.Sqrt(j)
		ins := fmt.Sprintf("%v", ms)
		if strings.Count(ins, ".") > 0 {
			continue
		}
		switch ms.(type) {
		default:
			continue

		case float64:
			//logger.Printf("this is float:", ms)
			//当j的平方根为 整数时

			//logger.Println("sqrt:", j, "result:", ins)
			sms := reverseTwo(ins)

			// 平方根的逆序 整数
			intsms, _ := strconv.ParseInt(sms, 10, 64)
			// j 的逆序 整数格式
			intsj, _ := strconv.ParseInt(sj, 10, 64)
			// 如果 j 的平方根 的逆序的 平方 等于 j 的逆序 则输出
			if intsms*intsms == intsj {
				logger.Println("strconv:", j, "sqrt:=", ms, "reverse j", intsj, "sqrt:=", intsms)
			}
			continue
		}
	}
}

func main() {
	runtime.ReadMemStats(&runtime.MemStats{})
	t0 := time.Now()
	//count(1000000000.0)
	logger.Println("cost time", time.Since(t0))

	c := func() {
		// Ask runtime.Callers for up to 10 PCs, including runtime.Callers itself.
		pc := make([]uintptr, 10)
		n := runtime.Callers(0, pc)
		if n == 0 {
			// No PCs available. This can happen if the first argument to
			// runtime.Callers is large.
			//
			// Return now to avoid processing the zero Frame that would
			// otherwise be returned by frames.Next below.
			return
		}

		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		// Loop to get frames.
		// A fixed number of PCs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()

			// Process this frame.
			//
			// To keep this example's output stable
			// even if there are changes in the testing package,
			// stop unwinding when we leave package runtime.
			if !strings.Contains(frame.File, "runtime/") {
				break
			}
			fmt.Printf("- more:%v | %s\n", more, frame.Function)

			// Check whether there are more frames to process after this one.
			if !more {
				break
			}
		}
	}

	b := func() { c() }
	a := func() { b() }

	a()
}
