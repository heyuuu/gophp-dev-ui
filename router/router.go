package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/controller"
	"net/http"
	_ "net/http/pprof"
	"path"
	"path/filepath"
)

func InitRouter(r *gin.Engine, staticFs embed.FS) {
	// static
	r.GET("/static/*filepath", func(c *gin.Context) {
		handleStatic(c, staticFs, c.Param("filepath"))
	})
	r.GET("/", controller.Index)

	// api
	apiGroup := r.Group("api")
	{
		apiGroup.POST("/run/code", controller.ApiHandler(controller.ApiRunCode))

		apiGroup.GET("/test/path_list", controller.ApiHandler(controller.TestPathList))
		apiGroup.GET("/test/list", controller.ApiHandler(controller.TestList))
		apiGroup.GET("/test/detail", controller.ApiHandler(controller.TestDetail))
		apiGroup.POST("/test/run", controller.ApiHandler(controller.TestRun))
		apiGroup.POST("/test/run_custom", controller.ApiHandler(controller.TestRunCustom))
	}

	// pprof
	r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
}

// 静态文件 MIME 类型映射
var staticContentTypes = map[string]string{
	".css":  "text/css; charset=utf-8",
	".js":   "text/javascript; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".png":  "image/png",
}

func handleStatic(c *gin.Context, staticFs embed.FS, pathInfo string) {
	if pathInfo == "" || pathInfo == "/" {
		pathInfo = "static/index.html"
	} else {
		pathInfo = path.Join("static", pathInfo)
	}

	// 尝试读取文件
	file, err := staticFs.ReadFile(pathInfo)
	if err != nil && pathInfo != "static/index.html" {
		pathInfo = "static/index.html"
		file, err = staticFs.ReadFile(pathInfo)
	}
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// 设置 Content-Type
	suffix := filepath.Ext(pathInfo)
	if contentType, ok := staticContentTypes[suffix]; ok {
		c.Header("Content-type", contentType)
	}

	// 返回文件内容
	c.Data(http.StatusOK, "", file)
}
