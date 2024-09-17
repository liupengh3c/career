package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func add(d int) {
	sum := 0
	for i := 0; i < d; i++ {
		sum += i
	}
	fmt.Println("the sum is: ", sum)
}
func main() {
	var wg sync.WaitGroup
	runTimes := 20
	now := time.Now()
	mpf, _ := ants.NewMultiPoolWithFunc(10, runTimes/10, func(i interface{}) {
		add(i.(int))
		wg.Done()
	}, ants.LeastTasks)
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		mpf.Invoke(1000000000)
	}
	wg.Wait()
	fmt.Println("循环运行：", time.Since(now))
}
