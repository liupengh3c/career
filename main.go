package main

import (
	"fmt"
	"math/rand"
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
	for i := 0; i < 100; i++ {
		fmt.Println(rand.Intn(100))
	}
}
