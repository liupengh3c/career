package main

import (
	"fmt"
	"sort"
)

type DwData struct {
	CarID           string `json:"car_id"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	MapRegion       string `json:"map_region"`
	HardwareVersion int    `json:"hardware_version"`
}

// 按起始时间排序的 ByStartTime 切片类型
type ByStartTime []DwData

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool { return a[i].StartTime < a[j].StartTime }

func main() {
	clips := []DwData{
		{
			CarID:           "1234567890",
			StartTime:       1200,
			EndTime:         1205,
			MapRegion:       "region_a",
			HardwareVersion: 1,
		},
		{
			CarID:           "1234567890",
			StartTime:       1100,
			EndTime:         1105,
			MapRegion:       "region_a",
			HardwareVersion: 1,
		},
		{
			CarID:           "1234567890",
			StartTime:       1206,
			EndTime:         1210,
			MapRegion:       "region_a",
			HardwareVersion: 1,
		},
		{
			CarID:           "1234567890",
			StartTime:       1106,
			EndTime:         1110,
			MapRegion:       "region_a",
			HardwareVersion: 1,
		},
		{
			CarID:           "1234567890",
			StartTime:       1306,
			EndTime:         1310,
			MapRegion:       "region_a",
			HardwareVersion: 1,
		},
	}

	sort.Sort(ByStartTime(clips))
	fmt.Println(clips) // 按起始时间排序后的切片,先打印出StartTime=1100的记录

}

func MergeClips(adcStatus []DwData) []DwData {
	adcTags := []DwData{}
	// 按起始时间排序
	sort.Sort(ByStartTime(adcStatus))
	if len(adcStatus) == 0 {
		return adcTags
	}
	// 初始化current变量，用于记录当前正在处理的片段
	current := adcStatus[0]
	for _, v := range adcStatus[1:] {
		if int64(current.EndTime*1000+1000) >= int64(v.StartTime*1000) {
			// 如果当前片段的结束时间大于等于下一个片段的开始时间，则合并这两个片段
			current.EndTime = v.EndTime
		} else {
			autoTag := DwData{
				CarID:     v.CarID,
				StartTime: int64(current.StartTime * 1000),
				EndTime:   int64(current.EndTime * 1000),
			}
			// 当前的时间片和上一个时间片不相邻，则将上一个片段加入到结果切片中
			adcTags = append(adcTags, autoTag)
			// 更新当前片段为下一个片段
			current = v
		}
	}
	autoTag := DwData{
		CarID:     current.CarID,
		StartTime: int64(current.StartTime * 1000),
		EndTime:   int64(current.EndTime * 1000),
	}
	adcTags = append(adcTags, autoTag)
	return adcTags
}
