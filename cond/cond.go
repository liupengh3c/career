package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func GeneralAttack(x int, num *int, wg *sync.WaitGroup, cond *sync.Cond) {
	defer wg.Done()
	// 随机一个10秒内的时间用于修整
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	fmt.Println("兵团", x+1, "在", time.Now().Format("2006-01-02 15:04:05"), "准备完毕，等待总攻")
	cond.L.Lock()
	// 准备完毕，等待教练发令
	cond.Wait()
	cond.L.Unlock()
	fmt.Println("兵团", x+1, "开始发起总攻", time.Now().Format("2006-01-02 15:04:05"))
}
func MainCond() {
	var wg sync.WaitGroup
	num := 0
	var lock sync.Mutex
	con := sync.NewCond(&lock)
	fmt.Println("将军下令，开始修整:", time.Now().Format("2006-01-02 15:04:05"), ",15s后开始总攻")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go GeneralAttack(i, &num, &wg, con)
	}
	// 等待发起总攻
	time.Sleep(15 * time.Second)
	fmt.Println("将军下令，开始总攻:", time.Now().Format("2006-01-02 15:04:05"))
	// 发起总攻
	con.Broadcast()
	wg.Wait()
}
