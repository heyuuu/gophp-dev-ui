package controller

import (
	"github.com/gin-gonic/gin"
	"gophp-dev-ui/pkg/service/runservice"
)

// api: /test/config
type testConfigParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
}

func TestConfigHandler(c *gin.Context) any {
	// 获取请求参数
	var p testConfigParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	return gin.H{
		"defaultTestRoot": manager.DefaultTestRoot(),
	}
}

// api: /test/path_list
type testPathListParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
	Root string             `form:"root" binding:"required"`
}

func TestPathListHandler(c *gin.Context) any {
	// 获取请求参数
	var p testPathListParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 获取目录
	list, err := manager.FindTestPaths(p.Root)
	if err != nil {
		return err
	}

	return gin.H{
		"root":  p.Root,
		"list":  list,
		"count": len(list),
	}
}

// api: /test/case_list
type testCaseListParam struct {
	Mode   runservice.RunMode `form:"mode" binding:"required"`
	Root   string             `form:"root" binding:"required"`
	Path   string             `form:"path"`
	Offset int                `form:"offset"`
	Limit  int                `form:"limit"`
}

func TestCaseListHandler(c *gin.Context) any {
	// 获取请求参数
	var p testCaseListParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 获取 case 列表
	list, err := manager.FindTestCases(p.Root, p.Path)
	if err != nil {
		return err
	}
	total := len(list)

	// offset && limit
	if p.Offset > 0 {
		list = list[min(p.Offset, total):]
	}
	if p.Limit > 0 && p.Limit < len(list) {
		list = list[:p.Limit]
	}

	return gin.H{
		"root":   p.Root,
		"path":   p.Path,
		"offset": p.Offset,
		"limit":  p.Limit,
		"list":   list,
		"total":  total,
		"count":  len(list),
	}
}

// api: /test/detail
type testDetailParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
	Root string             `form:"root" binding:"required"`
	Path string             `form:"path" binding:"required"`
}

func TestDetailHandler(c *gin.Context) any {
	// 获取请求参数
	var p testDetailParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 获取详情
	content, err := manager.TestCaseDetail(p.Root, p.Path)
	if err != nil {
		return err
	}

	return gin.H{
		"root":    p.Root,
		"path":    p.Path,
		"content": content,
	}
}

// api: /test/run
type testRunParam struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
	Root string             `form:"root" binding:"required"`
	Path string             `form:"path" binding:"required"`
}

func TestRunHandler(c *gin.Context) any {
	// 获取请求参数
	var p testRunParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 获取详情
	result, err := manager.RunTestCase(p.Root, p.Path)
	if err != nil {
		return err
	}

	return gin.H{
		"root":   p.Root,
		"path":   p.Path,
		"result": result,
	}
}

// api: /test/run_custom
type testRunCustomParam struct {
	Mode    runservice.RunMode `form:"mode" binding:"required"`
	Root    string             `form:"root" binding:"required"`
	Path    string             `form:"path" binding:"required"`
	Content string             `form:"content" binding:"required"`
}

func TestRunCustomHandler(c *gin.Context) any {
	// 获取请求参数
	var p testRunCustomParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 获取 manager
	manager, err := selectManager(p.Mode)
	if err != nil {
		return err
	}

	// 获取详情
	result, err := manager.RunTestCaseCustom(p.Root, p.Path, p.Content)
	if err != nil {
		return err
	}

	return gin.H{
		"root":   p.Root,
		"path":   p.Path,
		"result": result,
	}
}
