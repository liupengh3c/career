package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 常量定义
const (
	// 时间戳位数
	TimeBits = 28
	// 工作节点ID位数
	WorkerBits = 22
	// 序列号位数
	SeqBits = 13
	// 总位数
	TotalBits = 64
	// 自定义epoch: 使用固定时间点(2020-01-01 00:00:00 UTC)
	// FCoreEpoch = 1577836800
	FCoreEpoch = 1463673600
)

// BitsAllocator 位分配器
type BitsAllocator struct {
	SignBits      uint64
	TimestampBits uint64
	WorkerIdBits  uint64
	SequenceBits  uint64
	// 最大值
	MaxDeltaSeconds uint64
	MaxWorkerId     uint64
	MaxSequence     uint64
	// 位偏移
	TimestampShift uint64
	WorkerIdShift  uint64
	SequenceShift  uint64
}

// NewBitsAllocator 创建新的位分配器
func NewBitsAllocator(timeBits, workerBits, seqBits uint64) *BitsAllocator {
	signBits := uint64(1)
	totalBits := signBits + timeBits + workerBits + seqBits
	if totalBits != TotalBits {
		panic(fmt.Sprintf("Bits not equal to %d", TotalBits))
	}

	timestampShift := workerBits + seqBits
	workerIdShift := seqBits

	maxDeltaSeconds := uint64(1)<<timeBits - 1
	maxWorkerId := uint64(1)<<workerBits - 1
	maxSequence := uint64(1)<<seqBits - 1
	fmt.Println(maxDeltaSeconds, maxWorkerId, maxSequence)
	return &BitsAllocator{
		SignBits:        signBits,
		TimestampBits:   timeBits,
		WorkerIdBits:    workerBits,
		SequenceBits:    seqBits,
		MaxDeltaSeconds: maxDeltaSeconds,
		MaxWorkerId:     maxWorkerId,
		MaxSequence:     maxSequence,
		TimestampShift:  timestampShift,
		WorkerIdShift:   workerIdShift,
		SequenceShift:   0,
	}
}

// Allocate 分配UID
func (b *BitsAllocator) Allocate(deltaSeconds, workerId, sequence uint64) uint64 {
	uid := deltaSeconds << b.TimestampShift
	uid |= workerId << b.WorkerIdShift
	uid |= sequence
	return uid
}

// UidGenerator UID生成器接口
type UidGenerator interface {
	GetUID() (uint64, error)
	ParseUID(uid uint64) string
}

// DefaultUidGenerator 默认UID生成器
type DefaultUidGenerator struct {
	bitsAllocator *BitsAllocator
	workerId      uint64
	sequence      uint64
	lastSecond    int64
	mutex         sync.Mutex
}

// NewDefaultUidGenerator 创建默认UID生成器
func NewDefaultUidGenerator(workerId uint64) *DefaultUidGenerator {
	bitsAllocator := NewBitsAllocator(TimeBits, WorkerBits, SeqBits)
	return &DefaultUidGenerator{
		bitsAllocator: bitsAllocator,
		workerId:      workerId,
		sequence:      0,
		lastSecond:    -1,
	}
}

// GetUID 获取UID
func (g *DefaultUidGenerator) GetUID() (uint64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.nextId()
}

// getCurrentSecond 获取当前秒数
func (g *DefaultUidGenerator) getCurrentSecond() int64 {
	return time.Now().Unix()
}

// nextId 生成下一个ID
func (g *DefaultUidGenerator) nextId() (uint64, error) {
	currentSecond := g.getCurrentSecond()

	// 时钟回拨检查
	if currentSecond < g.lastSecond {
		refusedSeconds := g.lastSecond - currentSecond
		return 0, fmt.Errorf("Clock moved backwards. Refusing for %d seconds", refusedSeconds)
	}

	// 同一秒内增加序列号
	if currentSecond == g.lastSecond {
		g.sequence = (g.sequence + 1) & g.bitsAllocator.MaxSequence
		// 序列号耗尽，等待下一秒
		if g.sequence == 0 {
			currentSecond = g.getNextSecond(g.lastSecond)
		}
	} else {
		// 不同秒数，序列号重置
		g.sequence = 0
	}

	g.lastSecond = currentSecond

	// 计算deltaSeconds
	deltaSeconds := uint64(currentSecond - FCoreEpoch)

	// 检查时间戳是否耗尽
	if deltaSeconds > g.bitsAllocator.MaxDeltaSeconds {
		// 调整为取模运算，避免溢出
		deltaSeconds = deltaSeconds % (g.bitsAllocator.MaxDeltaSeconds + 1)
	}

	// 分配位
	return g.bitsAllocator.Allocate(deltaSeconds, g.workerId, g.sequence), nil
}

// getNextSecond 获取下一秒
func (g *DefaultUidGenerator) getNextSecond(lastTimestamp int64) int64 {
	timestamp := g.getCurrentSecond()
	for timestamp <= lastTimestamp {
		timestamp = g.getCurrentSecond()
	}
	return timestamp
}

// ParseUID 解析UID
func (g *DefaultUidGenerator) ParseUID(uid uint64) string {
	// 解析序列号
	sequence := (uid << (TotalBits - SeqBits)) >> (TotalBits - SeqBits)
	// 解析workerId
	workerId := (uid << (TimeBits + 1)) >> (TotalBits - WorkerBits)
	// 解析deltaSeconds
	deltaSeconds := uid >> (WorkerBits + SeqBits)

	// 计算实际时间
	actualSeconds := int64(FCoreEpoch) + int64(deltaSeconds)
	actualTime := time.Unix(actualSeconds, 0)

	return fmt.Sprintf("UID:%d parsed -> timestamp: %s, workerId: %d, sequence: %d",
		uid, actualTime.Format("2006-01-02 15:04:05"), workerId, sequence)
}

// CachedUidGenerator 缓存UID生成器
type CachedUidGenerator struct {
	*DefaultUidGenerator
	ringBuffer *RingBuffer
}

// RingBuffer 环形缓冲区
type RingBuffer struct {
	size       int
	buffer     []uint64
	readIndex  int
	writeIndex int
	mutex      sync.RWMutex
}

// NewRingBuffer 创建环形缓冲区
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		size:   size,
		buffer: make([]uint64, size),
	}
}

// take 从缓冲区取出UID
func (r *RingBuffer) take() (uint64, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.readIndex == r.writeIndex {
		// 实际应用中这里应该触发异步填充缓冲区
		return 0, errors.New("Buffer is empty, need padding")
	}

	uid := r.buffer[r.readIndex]
	r.readIndex = (r.readIndex + 1) % r.size
	return uid, nil
}

// put 向缓冲区添加UID
func (r *RingBuffer) put(uid uint64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	nextIndex := (r.writeIndex + 1) % r.size
	if nextIndex == r.readIndex {
		return errors.New("Buffer is full")
	}

	r.buffer[r.writeIndex] = uid
	r.writeIndex = nextIndex
	return nil
}

// 示例使用
func main() {
	// 测试默认生成器
	testDefaultGenerator()

	// 测试缓存生成器
	// testCachedGenerator()
}

func testDefaultGenerator() {
	gen := NewDefaultUidGenerator(1)

	// 生成10个UID测试
	for i := 0; i < 1; i++ {
		uid, err := gen.GetUID()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Default UID %d: %d\n", i+1, uid)
		fmt.Println(gen.ParseUID(uid))
	}
}

func testCachedGenerator() {
	gen := &CachedUidGenerator{
		DefaultUidGenerator: NewDefaultUidGenerator(1),
		ringBuffer:          NewRingBuffer(100),
	}

	// 预填充一些UID
	for i := 0; i < 1; i++ {
		uid, err := gen.DefaultUidGenerator.GetUID()
		if err != nil {
			fmt.Printf("Generate error: %v\n", err)
			continue
		}
		if err := gen.ringBuffer.put(uid); err != nil {
			fmt.Printf("Put to buffer error: %v\n", err)
		}

	}

	// 从缓存获取UID
	for i := 0; i < 1; i++ {
		uid, err := gen.GetUID()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Cached UID %d: %d\n", i+1, uid)
		fmt.Println(gen.ParseUID(uid))
	}
}
