package main

import (
	"fmt"
	"sync"
	"time"
)

// 与 Baidu UidGenerator 对齐的配置
const (
	TimeBits   = 28 // 时间位（秒级）
	WorkerBits = 22 // 机器位
	SeqBits    = 13 // 序列号位
	// EpochStart = 1577808000 // 2020-01-01 00:00:00 UTC （与FCoreEpoch一致）
	EpochStart = 1463673600 // 2016-05-20 00:00:00 UTC
)

var (
	MaxSequence = int64(1<<SeqBits - 1)
	MaxWorkerID = int64(1<<WorkerBits - 1)
)

type UidGenerator struct {
	workerID int64
	lastTime int64
	sequence int64
	mutex    sync.Mutex
}

func NewUidGenerator(workerID int64) *UidGenerator {
	if workerID > MaxWorkerID {
		panic("workerID too large")
	}
	return &UidGenerator{
		workerID: workerID,
	}
}

func (g *UidGenerator) NextUID() int64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	start := time.Now().Nanosecond()
	defer func() {
		fmt.Printf("耗时：%d 纳秒\n", time.Now().Nanosecond()-start)
	}()
	now := time.Now().Unix() - EpochStart

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & MaxSequence
		if g.sequence == 0 {
			for now <= g.lastTime {
				now = time.Now().Unix() - EpochStart
			}
		}
	} else {
		g.sequence = 0
	}
	g.lastTime = now
	// 注意：int64 左移不会丢符号位，这样可以产生负值
	uid := (now << (WorkerBits + SeqBits)) | (g.workerID << SeqBits) | g.sequence
	return uid
}

func main() {
	gen := NewUidGenerator(1)

	for i := 0; i < 10; i++ {
		fmt.Println(gen.NextUID())
	}
}
