package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heyuuu/gophp/tests"
	"os"
	"path/filepath"
)

type testPathListParam struct {
	Src string `form:"src" binding:"required"`
}

func TestPathList(c *gin.Context) ApiResult {
	var err error

	var p testPathListParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return apiError(err)
	}

	// testCases
	testPaths := tests.FindTestPathsInSrcDir(p.Src, true)

	return apiSucc(gin.H{
		"list":  testPaths,
		"count": len(testPaths),
	})
}

type testListParam struct {
	Src    string `form:"src" binding:"required"`
	Path   string `form:"path"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

func TestList(c *gin.Context) ApiResult {
	var err error

	var p testListParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return apiError(err)
	}

	var testCases []*tests.TestCase
	if p.Path == "" {
		testCases, err = tests.FindTestCasesInSrcDir(p.Src, false)
	} else {
		dir := filepath.Join(p.Src, p.Path)
		testCases, err = tests.FindTestCases(p.Src, dir)
	}
	if err != nil {
		return apiError(err)
	}

	// offset && limit
	total := len(testCases)
	if p.Offset > 0 {
		if len(testCases) > p.Offset {
			testCases = testCases[p.Offset:]
		} else {
			testCases = nil
		}
	}
	if p.Limit > 0 && p.Limit < len(testCases) {
		testCases = testCases[:p.Limit]
	}
	count := len(testCases)

	// fileNames
	var testNames = make([]string, len(testCases))
	for i, tc := range testCases {
		testNames[i] = tc.FileName()
	}

	return apiSucc(gin.H{
		"list":   testNames,
		"offset": p.Offset,
		"limit":  p.Limit,
		"total":  total,
		"count":  count,
	})
}

type testDetailParam struct {
	Src  string `form:"src" binding:"required"`
	Path string `form:"path" binding:"required"`
}

func TestDetail(c *gin.Context) ApiResult {
	var err error
	var p testDetailParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return apiError(err)
	}

	fileName := p.Path
	filePath := filepath.Join(p.Src, p.Path)

	tc := tests.NewTestCase(fileName, filePath)
	sections, err := tc.Parse()
	if err != nil {
		return apiError(fmt.Errorf("parse test-case file failed: %w", err))
	}

	return apiSucc(gin.H{
		"src":      p.Src,
		"path":     p.Path,
		"sections": sections,
	})
}

type testRunParams struct {
	Src  string `form:"src" binding:"required"`
	Path string `form:"path" binding:"required"`
}

func TestRun(c *gin.Context) ApiResult {
	var err error
	var p testRunParams
	if err = c.ShouldBindJSON(&p); err != nil {
		return apiError(err)
	}

	fileName := p.Path
	filePath := filepath.Join(p.Src, p.Path)

	tc := tests.NewTestCase(fileName, filePath)
	result := runTestCase(p.Src, tc)

	return apiSucc(gin.H{
		"fileName": tc.FileName(),
		"filePath": tc.FilePath(),
		"sections": tc.Sections(),
		"result":   result.MainType(),
		"output":   result.Output(),
	})
}

type testRunCustomParams struct {
	Src      string            `form:"src" binding:"required"`
	Path     string            `form:"path"`
	Sections map[string]string `form:"sections" binding:"required"`
}

func TestRunCustom(c *gin.Context) ApiResult {
	var err error
	var p testRunCustomParams
	if err = c.ShouldBindJSON(&p); err != nil {
		return apiError(err)
	}

	var fileName, filePath string
	if p.Path != "" {
		fileName = p.Path
		filePath = filepath.Join(p.Src, p.Path)
	} else {
		filePath, err = createTempTestFile()
		if err != nil {
			return apiError(err)
		}

		fileName = filepath.Base(filePath)
	}

	tc := tests.NewTestCaseParsed(fileName, filePath, p.Sections)
	result := runTestCase(p.Src, tc)

	return apiSucc(gin.H{
		"fileName": tc.FileName(),
		"filePath": tc.FilePath(),
		"sections": tc.Sections(),
		"result":   result.MainType(),
		"output":   result.Output(),
	})
}

func createTempTestFile() (string, error) {
	fs, err := os.CreateTemp(os.TempDir(), "gophp_dev_*.phpt")
	if err != nil {
		return "", err
	}
	defer fs.Close()

	return fs.Name(), nil
}

func runTestCase(src string, tc *tests.TestCase) *tests.Result {
	conf := tests.DefaultConfig()
	conf.SrcDir = src
	return tests.TestOneCase(conf, tc)
}
