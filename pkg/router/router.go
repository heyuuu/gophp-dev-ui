package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/pkg/controller"
	"net/http"
	_ "net/http/pprof"
	"path"
	"path/filepath"
)

func InitRouter(r *gin.Engine, staticFs embed.FS) {
	// 静态文件(指向 staticFs 对应文件，若文件不存在则指向 index.html)
	r.GET("/static/*filepath", func(c *gin.Context) {
		handleStatic(c, staticFs, c.Param("filepath"))
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static")
	})

	// api
	apiGroup := r.Group("api")
	{
		// run
		apiGroup.GET("/run/config", controller.ApiHandler(controller.RunConfigHandler))
		apiGroup.POST("/run/code", controller.ApiHandler(controller.RunCodeHandler))

		// test
		apiGroup.GET("/test/config", controller.ApiHandler(controller.TestConfigHandler))
		apiGroup.GET("/test/path_list", controller.ApiHandler(controller.TestPathListHandler))
		apiGroup.GET("/test/case_list", controller.ApiHandler(controller.TestCaseListHandler))
		apiGroup.GET("/test/detail", controller.ApiHandler(controller.TestDetailHandler))
		apiGroup.POST("/test/run", controller.ApiHandler(controller.TestRunHandler))
		apiGroup.POST("/test/run_custom", controller.ApiHandler(controller.TestRunCustomHandler))
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
