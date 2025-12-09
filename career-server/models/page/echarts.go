package page

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func CreateLineChart(ctx *gin.Context) *charts.Line {
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

	return line
}
func generateLineData() []opts.LineData {
	items := make([]opts.LineData, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		items = append(items, opts.LineData{Value: r.Intn(10000)})
	}
	return items
}
