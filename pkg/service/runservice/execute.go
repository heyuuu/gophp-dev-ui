package runservice

import (
	"github.com/heyuuu/gophp/kits/slicekit"
	"github.com/heyuuu/gophp/tests"
	"path/filepath"
)

type ExecuteManager struct {
}

var _ Manager = (*ExecuteManager)(nil)

func (e ExecuteManager) AllResultTypes() []ResultType {
	return []ResultType{
		RunTypeAst,
		RunTypeAstPrint,
		RunTypeExec,
		RunTypeExecRaw,
	}
}

func (e ExecuteManager) Run(code string) *RunResult {
	return new(executeRunner).run(code)
}

func (e ExecuteManager) DefaultTestRoot() string {
	return "/Users/heyu/Code/src/php-7.4.33"
}

func (e ExecuteManager) FindTestPaths(root string) ([]string, error) {
	return tests.FindTestPathsInSrcDir(root, true), nil
}

func (e ExecuteManager) FindTestCases(root string, path string) ([]string, error) {
	var err error

	var testCases []*tests.TestCase
	if path == "" {
		testCases, err = tests.FindTestCasesInSrcDir(root, false)
	} else {
		dir := filepath.Join(root, path)
		testCases, err = tests.FindTestCases(root, dir)
	}
	if err != nil {
		return nil, err
	}

	return slicekit.Map(testCases, func(t *tests.TestCase) string {
		return t.FileName()
	}), nil
}
