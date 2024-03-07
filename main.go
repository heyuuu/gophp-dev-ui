package main

import (
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	r.Run(":8080")
}
