package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "static")
}
