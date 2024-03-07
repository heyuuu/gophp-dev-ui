package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heyuuu/gophp/tests"
	"path/filepath"
)

type testsListParam struct {
	SrcDir string `form:"src" binding:"required"`
	Path   string `form:"path"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

func TestsList(c *gin.Context) ApiResult {
	var err error

	var p testsListParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return apiError(err)
	}

	var testNames []string
	if p.Path == "" {
		testNames, err = tests.FindTestFilesInSrcDir(p.SrcDir, false)
	} else {
		dir := filepath.Join(p.SrcDir, p.Path)
		testNames, err = tests.FindTestFiles(dir)
	}
	if err != nil {
		return apiError(err)
	}

	// offset && limit
	if p.Offset > 0 {
		if len(testNames) > p.Offset {
			testNames = testNames[p.Offset:]
		} else {
			testNames = nil
		}
	}
	total := len(testNames)
	if p.Limit > 0 && p.Limit < len(testNames) {
		testNames = testNames[:p.Limit]
	}
	for i, testFile := range testNames {
		testNames[i], _ = filepath.Rel(p.SrcDir, testFile)
	}
	count := len(testNames)

	return apiSucc(gin.H{
		"list":   testNames,
		"offset": p.Offset,
		"limit":  p.Limit,
		"total":  total,
		"count":  count,
	})
}

type testsDetailParam struct {
	SrcDir string `form:"src" binding:"required"`
	Path   string `form:"path" binding:"required"`
}

func TestsDetail(c *gin.Context) ApiResult {
	var err error
	var p testsDetailParam
	if err = c.ShouldBindQuery(&p); err != nil {
		return apiError(err)
	}

	shortFileName := p.Path
	realPath := filepath.Join(p.SrcDir, p.Path)

	tc, err := tests.ParseTestCase(0, realPath, shortFileName)
	if err != nil {
		return apiError(fmt.Errorf("parse test-case file failed: %w", err))
	}

	return apiSucc(gin.H{
		"sections": tc.Sections(),
	})
}

func TestsRun(c *gin.Context) ApiResult {
	return apiFail(0, "todo TestsRun")
}
