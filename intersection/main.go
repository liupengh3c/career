package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	UNKNOWN = iota
	A_CONTAINS_B
	A_OVERLAPS_LEFT_B
	A_OVERLAPS_RIGHT_B
	A_CONTAINS_BY_B
)

type ClipData struct {
	TaskId    string
	CarId     string
	StartTime int64
	EndTime   int64
	Expire    int64
	Topics    string
}

func IntervalRelation(clipA ClipData, clipB ClipData) int {
	// 情况 1: A 包含 B
	if clipA.StartTime <= clipB.StartTime && clipA.EndTime >= clipB.EndTime {
		return A_CONTAINS_B
	}

	// 情况 2: A 被包含于 B
	if clipB.StartTime < clipA.StartTime && clipB.EndTime > clipA.EndTime {
		return A_CONTAINS_BY_B
	}

	// 情况 3: A 左交 B
	if clipA.StartTime <= clipB.StartTime && (clipA.EndTime > clipB.StartTime && clipA.EndTime < clipB.EndTime) {
		return A_OVERLAPS_LEFT_B
	}

	// 情况 4: A 右交 B
	if clipA.EndTime >= clipB.EndTime && (clipA.StartTime < clipB.EndTime && clipA.StartTime > clipB.StartTime) {
		return A_OVERLAPS_RIGHT_B
	}
	return UNKNOWN
}

// func main() {
// 	clipA := ClipData{
// 		TaskId:    "1",
// 		CarId:     "10001",
// 		StartTime: 1672502000,
// 		EndTime:   1672505000,
// 		Expire:    1672506000,
// 		Topics:    "",
// 	}
// 	clipB := ClipData{
// 		TaskId:    "2",
// 		CarId:     "10001",
// 		StartTime: 1672502400,
// 		EndTime:   1672506000,
// 		Expire:    1672506000,
// 		Topics:    "",
// 	}
// 	res := IntervalRelation(clipA, clipB)
// 	switch res {
// 	case A_CONTAINS_B:
// 		println("A 包含 B")
// 	case A_CONTAINS_BY_B:
// 		println("A 被包含于 B")
// 	case A_OVERLAPS_LEFT_B:
// 		println("A 左交 B")
// 	case A_OVERLAPS_RIGHT_B:
// 		println("A 右交 B")
// 	default:
// 		println("未知关系")
// 	}
// }

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		fmt.Println("you input:", text) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
