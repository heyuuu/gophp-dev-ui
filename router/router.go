package router

import (
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/controller"
)

func InitRouter(r *gin.Engine) {
	// static
	r.Static("/static/", "static/")

	// api
	r.GET("/", controller.Index)
	r.GET("/tests/list", controller.ApiHandler(controller.TestsList))
	r.GET("/tests/detail", controller.ApiHandler(controller.TestsDetail))
	r.POST("/tests/run", controller.ApiHandler(controller.TestsRun))
}
