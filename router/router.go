package router

import (
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/controller"
	"os"
	"path/filepath"
)

func InitRouter(r *gin.Engine) {
	// static
	//r.Static("/static/asserts", "static/asserts")
	r.GET("/static/*filepath", func(c *gin.Context) {
		path := filepath.Join("static", c.Param("filepath"))
		if !isFileExist(path) {
			path = "static/index.html"
		}
		c.File(path)
	})
	r.GET("/", controller.Index)

	// api
	apiGroup := r.Group("api")
	{
		apiGroup.GET("/test/path_list", controller.ApiHandler(controller.TestPathList))
		apiGroup.GET("/test/list", controller.ApiHandler(controller.TestList))
		apiGroup.GET("/test/detail", controller.ApiHandler(controller.TestDetail))
		apiGroup.POST("/test/run_custom", controller.ApiHandler(controller.TestRunCustom))
	}
}

func isFileExist(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
