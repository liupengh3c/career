package main

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func main() {
	name := "/Users/liupeng/Documents/career/career/career-server/fit/1123-有氧.fit"
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	f, err := os.Open(name)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	defer f.Close()
	dec := decoder.New(f)
	fit, err := dec.Decode()
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	activity := filedef.NewActivity(fit.Messages...)
	fileId, _ := json.MarshalToString(activity.FileId)
	fmt.Println("activity", fileId)
	for k, v := range activity.Laps {
		speed := 1000.0 / (float32(v.AvgSpeed) / 1000.0)
		intSpeed := int(speed + 0.5)
		fmt.Printf("第%d公里, 配速:%v, 心率:%v\n", k+1, fmt.Sprintf("%v%02d", intSpeed/60, intSpeed%60), v.AvgHeartRate)
	}
	// fmt.Println("len record", len(activity.Records))
}
