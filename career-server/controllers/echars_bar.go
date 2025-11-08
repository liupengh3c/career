package controllers

import (
	"career-server/lib/result"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}
func barBasic() *charts.Bar {
	weeks := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic bar example", Subtitle: "This is the subtitle."}),
	)

	bar.SetXAxis(weeks).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())
	return bar
}
func CreateEchartsBar(ctx *gin.Context) {
	page := components.NewPage()
	response := &result.EchartsResponse{
		FilePath: "examples/bar.html",
	}

	bar := barBasic()
	f, err := os.Create(response.FilePath)
	if err != nil {
		fmt.Println("create file error:", err)
		return
	}
	defer f.Close()
	page.AddCharts(bar)
	page.Render(io.MultiWriter(f))
	response.EchoResult(ctx)
}
