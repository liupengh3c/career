package router

import (
	"career-server/controllers"
	"career-server/middlewares"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.Use(middlewares.Print, middlewares.GlobalLimiterMiddleware)
	// engine.Use(middlewares.Print, middlewares.IpLimitMiddleware)
	engine.POST("mianhuatang/image_search", controllers.ImageSearch)
	engine.POST("mianhuatang/load/articles", controllers.LoadArticles)
	engine.POST("mianhuatang/clip/search", controllers.ClipSearch)
}
