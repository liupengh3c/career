package router

import (
	"career-server/controllers"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	// engine.Use(middlewares.Print)
	engine.POST("mianhuatang/image_search", controllers.ImageSearch)
	engine.POST("mianhuatang/load/articles", controllers.LoadArticles)
}
