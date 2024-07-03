package runservice

import (
	"os"
	"path/filepath"
	"strings"
)

type CompileManager struct{}

var _ Manager = (*CompileManager)(nil)

func (c *CompileManager) AllResultTypes() []ResultType {
	return []ResultType{
		RunTypeSrc,
		RunTypeParseRaw,
		RunTypeAst,
		RunTypeAstPrint,
		RunTypeIr,
		RunTypeIrPrint,
		RunTypeIrRender,
		RunTypeExpected,
	}
}

func (c *CompileManager) Run(code string) *RunResult {
	//TODO implement me
	panic("implement me")
}

func (c *CompileManager) DefaultTestRoot() string {
	return "/Users/heyu/Code/sik/gophp-work-2/gophp/testdata"
}

func (c *CompileManager) FindTestPaths(root string) ([]string, error) {
	return findTestPaths(root, c.isTestCase), nil
}

func (c *CompileManager) FindTestCases(root string, path string) ([]string, error) {
	fullPath := filepath.Join(root, path)
	return findTestFiles(fullPath, c.isTestCase), nil
}

func (c *CompileManager) isTestCase(file string) bool {
	if !strings.HasSuffix(file, ".php") {
		return false
	}

	expectedFile := strings.TrimSuffix(file, ".php") + ".go"
	if _, err := os.Stat(expectedFile); err != nil {
		return false
	}

	return true
}
