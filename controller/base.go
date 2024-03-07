package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// -- result

type ApiResult struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  any    `json:"data"`
}

func apiSucc(data any) ApiResult {
	return ApiResult{Data: data}
}

func apiFail(code int, err string) ApiResult {
	return ApiResult{Code: code, Error: err}
}

func apiError(err error) ApiResult {
	var code int
	if coded, ok := err.(interface{ Code() int }); ok {
		code = coded.Code()
	}
	return ApiResult{Code: code, Error: err.Error()}
}

func ApiHandler(h func(c *gin.Context) ApiResult) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := h(c)
		c.JSON(http.StatusOK, r)
		return
	}
}

// -- params

func ParseQuery[T any](c *gin.Context) (T, error) {
	var tmp T
	err := c.ShouldBindQuery(&tmp)
	return tmp, err
}
