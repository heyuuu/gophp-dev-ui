package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heyuuu/gophp/tests"
	"gophp-dev-ui/pkg/service/runservice"
	"os"
	"path/filepath"
	"strings"
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
	var err error
	var p testDetailParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return err
	}

	fileName := p.Path
	filePath := filepath.Join(p.Root, p.Path)

	tc := tests.NewTestCase(fileName, filePath)
	sections, err := tc.Parse()
	if err != nil {
		return fmt.Errorf("parse test-case file failed: %w", err)
	}

	return gin.H{
		"root":     p.Root,
		"path":     p.Path,
		"sections": sections,
	}
}

// api: /test/run
type testRunParams struct {
	Mode runservice.RunMode `form:"mode" binding:"required"`
	Root string             `form:"root" binding:"required"`
	Path string             `form:"path" binding:"required"`
}

func TestRunHandler(c *gin.Context) any {
	var err error
	var p testRunParams
	if err = c.ShouldBindJSON(&p); err != nil {
		return err
	}

	fileName := p.Path
	filePath := filepath.Join(p.Root, p.Path)

	tc := tests.NewTestCase(fileName, filePath)
	ret := runTestCaseAndReturn(p.Root, tc)
	return ret
}

// api: /test/run_custom
type testRunCustomParams struct {
	Mode     runservice.RunMode `form:"mode" binding:"required"`
	Root     string             `form:"root" binding:"required"`
	Path     string             `form:"path" binding:"required"`
	Sections map[string]string  `form:"sections" binding:"required"`
}

func TestRunCustomHandler(c *gin.Context) any {
	var err error
	var p testRunCustomParams
	if err = c.ShouldBindJSON(&p); err != nil {
		return err
	}

	var fileName, filePath string
	if p.Path != "" {
		fileName = p.Path
		filePath = filepath.Join(p.Root, p.Path)
	} else {
		filePath, err = createTempTestFile()
		if err != nil {
			return err
		}

		fileName = filepath.Base(filePath)
	}

	tc := tests.NewTestCaseParsed(fileName, filePath, p.Sections)
	ret := runTestCaseAndReturn(p.Root, tc)
	return ret
}

func createTempTestFile() (string, error) {
	fs, err := os.CreateTemp(os.TempDir(), "gophp_dev_*.phpt")
	if err != nil {
		return "", err
	}
	defer fs.Close()

	return fs.Name(), nil
}

func runTestCase(src string, tc *tests.TestCase) (result *tests.Result, log string) {
	conf := tests.DefaultConfig()
	conf.SrcDir = src

	var buf strings.Builder
	conf.Logger = tests.LoggerFunc(func(tc *tests.TestCase, event int, message string) {
		if tc != nil {
			buf.WriteString(message)
		}
	})

	return tests.TestOneCase(conf, tc), buf.String()
}

func runTestCaseAndReturn(src string, tc *tests.TestCase) gin.H {
	result, log := runTestCase(src, tc)
	sections := tc.Sections()
	return gin.H{
		"fileName": tc.FileName(),
		"filePath": tc.FilePath(),

		// case
		//"sections": tc.Sections(),
		"code":   sections["FILE"],
		"expect": sections["EXPECT"] + sections["EXPECTF"] + sections["EXPECTREGEX"],

		// result
		"status":  result.MainType(),
		"output":  result.Output(),
		"info":    result.Info() + "\n" + log,
		"useTime": result.UseTime().Nanoseconds(),
	}
}
