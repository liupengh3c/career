package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Clip struct {
	TaskId     string
	CarId      string
	StartTime  int64
	EndTime    int64
	DeleteTime int64 // 数据删除时间
	Topics     string
}

// mergeBaseWithClips 将一个基准时间片与 n 个无交集时间片 merge
func mergeBaseWithClips(base Clip, others []Clip) []Clip {
	timePointsAll := []int64{base.StartTime, base.EndTime}
	for _, c := range others {
		timePointsAll = append(timePointsAll, c.StartTime, c.EndTime)
	}
	// 排序并去重
	// sort.Slice(timePointsAll, func(i, j int) bool { return timePointsAll[i] < timePointsAll[j] })
	slices.Sort(timePointsAll)
	uniq := []int64{timePointsAll[0]}
	for _, t := range timePointsAll[1:] {
		if t != uniq[len(uniq)-1] {
			uniq = append(uniq, t)
		}
	}

	mergeDatas := []Clip{}
	for i := 0; i < len(uniq)-1; i++ {
		clip := Clip{}
		s := uniq[i]
		e := uniq[i+1]
		if s >= e {
			continue
		}
		topics := ""
		expire := base.DeleteTime

		covered := false
		// 判断 base 覆盖
		if (base.StartTime <= s) && (base.EndTime >= e) {
			topicSli := strings.Split(base.Topics, ",")
			fmt.Println("topicSli:", topicSli)
			slices.Sort(topicSli)
			topics = strings.Join(topicSli, ",")
			fmt.Println("topicSli sorted:", topics)
			covered = true
			clip = base
			clip.StartTime = s
			clip.EndTime = e
			clip.Topics = topics
		}

		// 判断其他 clip 覆盖
		for _, c := range others {
			if (c.StartTime <= s) && (c.EndTime >= e) {
				clip = c
				clip.StartTime = s
				clip.EndTime = e
				// topics 处理逻辑：合并
				clip.Topics = mergeTopics(topics, c.Topics)
				if expire > c.DeleteTime {
					clip.DeleteTime = expire
				}
				covered = true
			}
		}

		if covered {
			mergeDatas = append(mergeDatas, clip)
		}
	}
	fmt.Println("segments:", mergeDatas)
	return mergeAdjacent(mergeDatas)
}

func mergeTopics(t1, t2 string) string {
	seen := make(map[string]struct{})
	result := []string{}
	t1s := strings.Split(t1, ",")
	t2s := strings.Split(t2, ",")
	for _, t := range t1s {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			result = append(result, t)
		}
	}
	for _, t := range t2s {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			result = append(result, t)
		}
	}
	slices.Sort(result)
	return strings.Join(result, ",")
}

// mergeAdjacent 合并相邻且 label/expire 相同的 segment
func mergeAdjacent(segs []Clip) []Clip {
	if len(segs) == 0 {
		return segs
	}
	res := []Clip{segs[0]}
	for i := 1; i < len(segs); i++ {
		last := &res[len(res)-1]
		cur := segs[i]
		if last.Topics == cur.Topics && last.DeleteTime == cur.DeleteTime {
			// 合并：延长 last 的 End
			last.EndTime = cur.EndTime
		} else {
			res = append(res, cur)
		}
	}
	return res
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

func main() {
	base := Clip{
		StartTime:  parseDate("2025-10-01 00:00:00").Unix(),
		EndTime:    parseDate("2025-10-20 00:00:00").Unix(),
		DeleteTime: parseDate("2025-11-02 00:00:00").Unix(),
		Topics:     "B,A",
	}

	others := []Clip{
		{
			StartTime:  parseDate("2025-10-05 00:00:00").Unix(),
			EndTime:    parseDate("2025-10-08 00:00:00").Unix(),
			DeleteTime: parseDate("2025-11-01 00:00:00").Unix(), // 与 base 相同
			Topics:     "A,B",
		},
		{
			StartTime:  parseDate("2025-10-12 00:00:00").Unix(),
			EndTime:    parseDate("2025-10-15 00:00:00").Unix(),
			DeleteTime: parseDate("2025-10-20 00:00:00").Unix(),
			Topics:     "B,C",
		},
		{
			StartTime:  parseDate("2025-10-02 00:00:00").Unix(),
			EndTime:    parseDate("2025-10-04 00:00:00").Unix(),
			DeleteTime: parseDate("2025-10-20 00:00:00").Unix(),
			Topics:     "B,C,D",
		},
	}

	segs := mergeBaseWithClips(base, others)
	for i, seg := range segs {
		fmt.Printf("[%d] %s -> %s | Label=%s | Expire=%s\n",
			i+1,
			time.Unix(seg.StartTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(seg.EndTime, 0).Format("2006-01-02 15:04:05"),
			seg.Topics,
			time.Unix(seg.DeleteTime, 0).Format("2006-01-02 15:04:05"),
		)
	}
}
