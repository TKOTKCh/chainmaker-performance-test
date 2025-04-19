package main

import (
	"fmt"
	"time"
)

//go:wasmexport sum
func sum(x, y int64) int64 {
	return x + y
}

//go:wasmexport Time2Str
func Time2Str() {
	aTime := time.Now()
	pattern := "2006-01-02 15:04:05.000"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(aTime.In(loc).Format(pattern))

}

func main() {

}
