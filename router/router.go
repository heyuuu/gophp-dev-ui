package router

import (
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/controller"
)

func InitRouter(r *gin.Engine) {
	// static
	r.Static("/static/", "static/")
	r.GET("/", controller.Index)

	// api
	r.GET("/test/path_list", controller.ApiHandler(controller.TestPathList))
	r.GET("/test/list", controller.ApiHandler(controller.TestList))
	r.GET("/test/detail", controller.ApiHandler(controller.TestDetail))
	r.POST("/test/run_custom", controller.ApiHandler(controller.TestRunCustom))
}
