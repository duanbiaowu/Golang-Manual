package structural

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Operator func(int) int

func OpDecorator(fn Operator) Operator {
	return func(n int) int {
		result := fn(n)
		return result
	}
}

type SumFunc func(int64, int64) int64

func getFuncName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func timedSumDecorator(f SumFunc) SumFunc {
	return func(start int64, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %d\n", getFuncName(f), time.Since(t))
		}(time.Now())
		return f(start, end)
	}
}
