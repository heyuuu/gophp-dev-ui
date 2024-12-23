package runservice

import (
	"github.com/heyuuu/gophp/compile/parser"
	"github.com/heyuuu/gophp/compile/render"
	"github.com/heyuuu/gophp/compile/transformer"
	"github.com/heyuuu/gophp/kits/oskit"
	"os"
	"path/filepath"
	"strings"
)

type CompileCase struct {
	src      string
	expected string
}

type CompileManager struct{}

var _ Manager = (*CompileManager)(nil)

func (m *CompileManager) AllResultTypes() []ResultType {
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

func (m *CompileManager) RunCode(code string) *RunResult {
	return RunCode(code, m.AllResultTypes())
}

func (m *CompileManager) DefaultTestRoot() string {
	return "/Users/heyu/Code/sik/gophp-work-2/gophp/testdata"
}

func (m *CompileManager) FindTestPaths(root string) ([]string, error) {
	return findTestPaths(root, m.isTestCase), nil
}

func (m *CompileManager) FindTestCases(root string, path string) ([]string, error) {
	return findTestFiles(root, path, m.isTestCase), nil
}

func (m *CompileManager) isTestCase(file string) bool {
	if !strings.HasSuffix(file, ".php") {
		return false
	}

	expectedFile := strings.TrimSuffix(file, ".php") + ".go"
	if _, err := os.Stat(expectedFile); err != nil {
		return false
	}

	return true
}

func (m *CompileManager) TestCaseDetail(root string, path string) (string, error) {
	fullPath := filepath.Join(root, path)
	return oskit.ReadFileAsString(fullPath)
}

func (m *CompileManager) RunTestCase(root string, path string) (*TestResult, error) {
	fullPath := filepath.Join(root, path)
	tc, err := m.loadTestCase(fullPath)
	if err != nil {
		return nil, err
	}

	return m.runTestCase(tc)
}

func (m *CompileManager) RunTestCaseCustom(root string, path string, content string) (*TestResult, error) {
	fullPath := filepath.Join(root, path)
	tc, err := m.loadTestCase(fullPath)
	if err != nil {
		return nil, err
	}
	tc.src = content

	return m.runTestCase(tc)
}

func (m *CompileManager) loadTestCase(srcFile string) (*CompileCase, error) {
	src, expected, err := loadTestCase(srcFile)
	if err != nil {
		return nil, err
	}
	tc := &CompileCase{
		src:      src,
		expected: expected,
	}
	return tc, nil
}

func (m *CompileManager) runTestCase(tc *CompileCase) (*TestResult, error) {
	output, err := m.runTestCode(tc.src)
	if err != nil {
		return nil, err
	}

	var status TestResultStatus
	if output == tc.expected {
		status = TestResultPass
	} else {
		status = TestResultFail
	}

	result := &TestResult{
		Code:     tc.src,
		Expected: tc.expected,
		Status:   status,
		Output:   output,
	}
	return result, nil
}

func (m *CompileManager) runTestCode(code string) (string, error) {
	astFile, err := parser.ParseCode(code)
	if err != nil {
		return "", err
	}

	irFile, err := transformer.TransformFile(astFile)
	if err != nil {
		return "", err
	}

	irRender, err := render.RenderFile(irFile)
	if err != nil {
		return "", err
	}

	return irRender, nil
}
