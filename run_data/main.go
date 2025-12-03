package main

import (
	"fmt"
	"os"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func main() {
	f, err := os.Open("/Users/liupeng/Documents/career/career/run_data/473297413331779787.fit")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := decoder.New(f)

	fit, err := dec.Decode()
	if err != nil {
		panic(err)
	}

	activity := filedef.NewActivity(fit.Messages...)

	fmt.Printf("File Type: %s\n", activity.FileId.Type)
	fmt.Printf("Sessions count: %d\n", len(activity.Sessions))
	fmt.Printf("Laps count: %d\n", len(activity.Laps))
	fmt.Printf("Records count: %d\n", len(activity.Records))

	for index, v := range activity.Laps {
		fmt.Printf("==========第%v公里=========\n", index+1)
		// fmt.Println("距离", v.TotalDistance/100000, "km")
		fmt.Println("平均心率", v.AvgHeartRate)
		fmt.Println("平均速度", v.AvgSpeed)
		fmt.Println("平均配速", Convert(float64(v.AvgSpeed)/1000.0))
		fmt.Println("运动时间", v.TotalElapsedTime)
	}
	// for _, v := range activity.Records {
	// 	fmt.Printf("  Distance: %g m\n", v.DistanceScaled())
	// 	fmt.Printf("  Lat: %d semicircles\n", v.PositionLat)
	// 	fmt.Printf("  Long: %d semicircles\n", v.PositionLong)
	// 	fmt.Printf("  Speed: %g m/s\n", v.SpeedScaled())
	// 	fmt.Printf("  HeartRate: %d bpm\n", v.HeartRate)
	// 	fmt.Printf("  Cadence: %d rpm\n", v.Cadence)
	// }

	// Output:
	// File Type: activity
	// Sessions count: 1
	// Laps count: 1
	// Records count: 3601
	//
	// Sample value of record[100]:
	//   Distance: 100 m
	//   Lat: 0 semicircles
	//   Long: 10717 semicircles
	//   Speed: 1 m/s
	//   HeartRate: 126 bpm
	//   Cadence: 100 rpm
}

func Convert(speed float64) string {
	dur := 1000.0/speed + 0.5
	m := int(dur) / 60
	s := int(dur) % 60
	return fmt.Sprintf("%v'%02d\"", m, s)
}
