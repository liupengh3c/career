package controllers

import (
	"career-server/models/page"

	"github.com/gin-gonic/gin"
)

func CreateEchartsLines(ctx *gin.Context) {
	// f, _ := os.Create("lines.html")
	// defer f.Close()
	line := page.CreateLineChart(ctx)
	// line.Render(f)
	// ctx.File(f.Name())
	line.Render(ctx.Writer)
}
