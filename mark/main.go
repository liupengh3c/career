package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Clip struct {
	StartTime int64
	EndTime   int64
	Expire    int64
	Label     string
	ID        string
}

type Segment struct {
	StartTime int64
	EndTime   int64
	Label     string
	Expire    int64
}

// mergeBaseWithClips 将一个基准时间片与 n 个无交集时间片 merge
func mergeBaseWithClips(base Clip, others []Clip) []Segment {
	points := []int64{base.StartTime, base.EndTime}
	for _, c := range others {
		points = append(points, c.StartTime, c.EndTime)
	}
	// 排序并去重
	sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
	uniq := []int64{points[0]}
	for _, t := range points[1:] {
		if t != uniq[len(uniq)-1] {
			uniq = append(uniq, t)
		}
	}

	segments := []Segment{}
	for i := 0; i < len(uniq)-1; i++ {
		s := uniq[i]
		e := uniq[i+1]
		if s >= e {
			continue
		}
		label := ""
		expire := base.Expire

		covered := false
		// 判断 base 覆盖
		if (base.StartTime <= s) &&
			(base.EndTime >= e) {
			label = base.Label
			covered = true
		}

		// 判断其他 clip 覆盖
		for _, c := range others {
			if (c.StartTime <= s) &&
				(c.EndTime >= e) {
				if !strings.Contains(label, c.Label) {
					label = label + c.Label
				}
				if c.Expire > expire {
					expire = c.Expire
				}
				covered = true
			}
		}

		if covered {
			segments = append(segments, Segment{
				StartTime: s,
				EndTime:   e,
				Label:     label,
				Expire:    expire,
			})
		}
	}
	fmt.Println("segments:", segments)
	return mergeAdjacent(segments)
}

// mergeAdjacent 合并相邻且 label/expire 相同的 segment
func mergeAdjacent(segs []Segment) []Segment {
	if len(segs) == 0 {
		return segs
	}
	res := []Segment{segs[0]}
	for i := 1; i < len(segs); i++ {
		last := &res[len(res)-1]
		cur := segs[i]
		if last.Label == cur.Label && last.Expire == cur.Expire {
			// 合并：延长 last 的 End
			last.EndTime = cur.EndTime
		} else {
			res = append(res, cur)
		}
	}
	return res
}

func parseDate(s string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return t
}

func main() {
	base := Clip{
		StartTime: parseDate("2025-10-05 00:00:00").Unix(),
		EndTime:   parseDate("2025-10-20 00:00:00").Unix(),
		Expire:    parseDate("2025-11-01 00:00:00").Unix(),
		Label:     "A",
		ID:        "base",
	}

	others := []Clip{
		{
			StartTime: parseDate("2025-10-01 00:00:00").Unix(),
			EndTime:   parseDate("2025-10-08 00:00:00").Unix(),
			Expire:    parseDate("2025-12-01 00:00:00").Unix(),
			Label:     "A",
			ID:        "c1",
		},
		{
			StartTime: parseDate("2025-10-12 00:00:00").Unix(),
			EndTime:   parseDate("2025-10-15 00:00:00").Unix(),
			Expire:    parseDate("2025-10-20 00:00:00").Unix(),
			Label:     "B",
			ID:        "c2",
		},
	}

	segs := mergeBaseWithClips(base, others)
	for i, seg := range segs {
		fmt.Printf("[%d] %s -> %s | Label=%s | Expire=%s\n",
			i+1,
			time.Unix(seg.StartTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(seg.EndTime, 0).Format("2006-01-02 15:04:05"),
			seg.Label,
			time.Unix(seg.Expire, 0).Format("2006-01-02 15:04:05"),
		)
	}
	// today, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	// n := today.Unix() * 1000
	// m := today.Add(24*time.Hour).Unix()*1000 - 1
	// fmt.Println(n, m)
	// 使用 UTC 确保不受时区影响
	now := time.Now().UTC()

	oneMonthAgo := now.AddDate(0, -1, 0)
	twoMonthsAgo := now.AddDate(0, -2, 0)

	fmt.Println("Now (UTC):        ", now.Format("2006-01-02 15:04:05"))
	fmt.Println("1 Month Ago (UTC):", oneMonthAgo.Format("2006-01-02 15:04:05"))
	fmt.Println("2 Months Ago (UTC):", twoMonthsAgo.Format("2006-01-02 15:04:05"))
}
