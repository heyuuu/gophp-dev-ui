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

func (e ExecuteManager) FindTestPaths(root string) []string {
	return tests.FindTestPathsInSrcDir(root, true)
}

func (e ExecuteManager) FindTestCases(path string) []string {
	var err error

	src, path := path, ""

	var testCases []*tests.TestCase
	if path == "" {
		testCases, err = tests.FindTestCasesInSrcDir(src, false)
	} else {
		dir := filepath.Join(src, path)
		testCases, err = tests.FindTestCases(src, dir)
	}
	if err != nil {
		panic(err)
	}

	return slicekit.Map(testCases, func(t *tests.TestCase) string {
		return t.FileName()
	})
}
