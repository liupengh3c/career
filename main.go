package main

import (
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"time"
)

type Person struct {
	Name  string `json:"name"`
	Age   *int32 `json:"age"`
	Count *int32 `json:"count"`
}

func printMemInfo() {
	for {
		m := runtime.MemStats{}
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
		fmt.Printf("\tTotal Alloc = %v MiB", m.TotalAlloc/1024/1024)
		fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
		fmt.Printf("\tNumGC = %v GCs\n", m.NumGC)
		time.Sleep(2 * time.Second)
	}
}
func main() {
	a := map[string]int32{}
	a["a"] = 1
	a["b"] = 2
	fmt.Println(a)
}
