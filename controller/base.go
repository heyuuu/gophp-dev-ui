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

const (
	CodeSucc = 0
	CodeFail = 1
)

func apiSucc(data any) ApiResult {
	return ApiResult{Code: CodeSucc, Data: data}
}

func apiFail(code int, err string) ApiResult {
	if code == CodeSucc {
		code = CodeFail
	}
	return ApiResult{Code: code, Error: err}
}

func apiError(err error) ApiResult {
	var code int
	if coded, ok := err.(interface{ Code() int }); ok {
		code = coded.Code()
	}
	return apiFail(code, err.Error())
}

func wrapApiResult(v any) ApiResult {
	switch value := v.(type) {
	case ApiResult:
		return value
	case error:
		return apiError(value)
	default:
		return apiSucc(value)
	}
}

func ApiHandler(h func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := wrapApiResult(h(c))
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
