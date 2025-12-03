package controllers

import (
	"career-server/models/page"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 渲染图表到HTML
func renderChart(c *gin.Context, chart interface{}) {
	err := chart.(interface {
		Render(w io.Writer) error
	}).Render(c.Writer)

	if err != nil {
		c.String(http.StatusInternalServerError, "图表渲染失败: "+err.Error())
		return
	}
}

func CreateEcharts(ctx *gin.Context) {
	bar := page.CreateLines(ctx)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	bar.Render(ctx.Writer)
}
