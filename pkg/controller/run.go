package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/pkg/service/runservice"
)

func selectManager(mode runservice.RunMode) (runservice.Manager, error) {
	manager := runservice.GetManager(mode)
	if manager == nil {
		return nil, errors.New("不支持此mode: " + string(mode))
	}

	return manager, nil
}

// api: /run/config
type runConfigParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
}

func RunConfigHandler(c *gin.Context) any {
	// 获取请求参数
	var p runConfigParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	return gin.H{
		"types": manager.AllResultTypes(),
	}
}

// api: /run/code
type runCodeParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
	Code string             `form:"code" binding:"required"`
}

func RunCodeHandler(c *gin.Context) any {
	// 获取请求参数
	var p runCodeParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 执行代码，并返回结果
	return manager.RunCode(p.Code)
}
