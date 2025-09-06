package main

import "fmt"

type MyDefine interface {
	~int | ~int32 | ~string
}

func SwapInt(a, b *int) {
	tmp := *a
	*a = *b
	*b = tmp
}
func SwapFloat(a, b *float32) {
	tmp := *a
	*a = *b
	*b = tmp
}
func SwapString(a, b *string) {
	tmp := *a
	*a = *b
	*b = tmp
}
func Swap[T any](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

func SwapConstrained[T MyDefine](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

type MyInt int

// 使用示例
func main() {
	c := 1
	d := 2
	SwapConstrained(&c, &d)
	fmt.Println(c, d)
	var g, h MyInt = 10, 20
	SwapConstrained(&g, &h)
}
