package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: r.Intn(100)})
	}
	return items
}

func barBasic() *charts.Bar {
	bar := charts.NewBar()
	weeks := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "每周销量统计", Subtitle: "苹果 vs 香蕉."}),
	)

	bar.SetXAxis(weeks).
		AddSeries("苹果", generateBarItems()).
		AddSeries("香蕉", generateBarItems())
	return bar
}

func CreateBar() {
	bar := barBasic()
	f, _ := os.Create("bar.html")
	defer f.Close()
	bar.Render(f)
}

func generateLineData() []opts.LineData {
	items := make([]opts.LineData, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		items = append(items, opts.LineData{Value: r.Intn(10000)})
	}
	return items
}
func CreateLineChart() {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "小度眼镜 vs 小度智能音箱",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "dark", // 使用暗黑主题
		}),
	)

	line.SetXAxis([]string{"1月", "2月", "3月", "4月", "5月", "6月"}).
		AddSeries("小度眼镜", generateLineData()).
		AddSeries("小度智能音箱", generateLineData()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	f, _ := os.Create("line.html")
	line.Render(f)
}

func CreatePieChart() {
	pie := charts.NewPie()

	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "用户设备分布",
		}),
	)

	pie.AddSeries("设备", []opts.PieData{
		{Name: "手机", Value: 1230},
		{Name: "平板", Value: 456},
		{Name: "电脑", Value: 789},
		{Name: "其他", Value: 123},
	})

	f, _ := os.Create("pie.html")
	pie.Render(f)
}

func main() {
	// CreateBar()
	CreateLineChart()
	// CreatePieChart()
}
