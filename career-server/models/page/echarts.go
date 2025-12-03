package page

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func GetRunDataFromFit(name string) ([]string, []opts.LineData) {
	fmt.Println("file_name", name)
	items := make([]opts.LineData, 0)
	x := []string{}
	f, err := os.Open(name)
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

	for index, v := range activity.Laps {
		x = append(x, fmt.Sprintf("%vkm", index+1))
		items = append(items, opts.LineData{Value: float32(v.AvgSpeed/100) / float32(v.AvgHeartRate)})
	}
	return x, items
}

func ScanAllFits(filePath string) []string {
	files := []string{}
	entrys, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Println("read dir error", err.Error())
		return files
	}
	for _, entry := range entrys {
		if entry.IsDir() || entry.Name() == ".DS_Store" {
			continue
		}
		files = append(files, filePath+entry.Name())
	}
	return files
}

func GetAllFitsData() ([]string, [][]string, [][]opts.LineData) {
	xAxis := [][]string{}
	yDatas := [][]opts.LineData{}
	base := "/Users/liupeng/Documents/career/career/career-server/fit/"
	fits := ScanAllFits(base)
	for _, fit := range fits {
		x, y := GetRunDataFromFit(fit)
		xAxis = append(xAxis, x)
		yDatas = append(yDatas, y)
	}
	return fits, xAxis, yDatas
}

func GetAllFitsName() []string {
	names := []string{}
	base := "/Users/liupeng/Documents/career/career/career-server/fit/"
	entrys, _ := os.ReadDir(base)
	for _, v := range entrys {
		if v.IsDir() {
			continue
		}
		names = append(names, strings.Split(v.Name(), ".")[0])
	}
	return names
}
func CreateLines(ctx *gin.Context) *charts.Line {
	fileNames, xs, datas := GetAllFitsData()
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "配速/心率",
		}),
	)

	fmt.Println("all names", fileNames)
	fmt.Println("all names", len(datas))
	for k, d := range datas {
		// fmt.Println(k, xs[k], fileNames[k])
		t := strings.Split(fileNames[k], "/")
		name := t[len(t)-1]
		line.SetXAxis(xs[k]).AddSeries(strings.Split(name, ".")[0], d)
	}
	line.SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: opts.Bool(true),
		}))
	return line
}
