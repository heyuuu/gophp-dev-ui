package runservice

import (
	"github.com/heyuuu/gophp/kits/oskit"
	"github.com/heyuuu/gophp/kits/slicekit"
	"github.com/heyuuu/gophp/tests"
	"os"
	"path/filepath"
	"strings"
)

type ExecuteManager struct {
}

var _ Manager = (*ExecuteManager)(nil)

func (m *ExecuteManager) AllResultTypes() []ResultType {
	return []ResultType{
		RunTypeAst,
		RunTypeAstPrint,
		RunTypeExec,
		RunTypeExecRaw,
	}
}

func (m *ExecuteManager) RunCode(code string) *RunResult {
	return new(executeRunner).run(code)
}

func (m *ExecuteManager) DefaultTestRoot() string {
	return "/Users/heyu/Code/src/php-7.4.33"
}

func (m *ExecuteManager) FindTestPaths(root string) ([]string, error) {
	return tests.FindTestPathsInSrcDir(root, true), nil
}

func (m *ExecuteManager) FindTestCases(root string, path string) ([]string, error) {
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

func (m *ExecuteManager) TestCaseDetail(root string, path string) (string, error) {
	//fileName := p.Path
	//filePath := filepath.Join(p.Root, p.Path)
	//
	//tc := tests.NewTestCase(fileName, filePath)
	//sections, err := tc.Parse()
	//if err != nil {
	//	return fmt.Errorf("parse test-case file failed: %w", err)
	//}

	fullPath := filepath.Join(root, path)
	return oskit.ReadFileAsString(fullPath)
}

func (m *ExecuteManager) RunTestCase(root string, path string) (*TestResult, error) {
	fileName := path
	filePath := filepath.Join(root, path)

	tc := tests.NewTestCase(fileName, filePath)
	return m.runTestCaseAndReturn(root, tc), nil
}

func (m *ExecuteManager) RunTestCaseCustom(root string, path string, content string) (*TestResult, error) {

	var fileName, filePath string
	var err error
	if path != "" {
		fileName = path
		filePath = filepath.Join(root, path)
	} else {
		filePath, err = m.createTempTestFile()
		if err != nil {
			return nil, err
		}

		fileName = filepath.Base(filePath)
	}

	tc := tests.NewTestCaseCustom(fileName, filePath, content)
	return m.runTestCaseAndReturn(root, tc), nil
}

func (m *ExecuteManager) createTempTestFile() (string, error) {
	fs, err := os.CreateTemp(os.TempDir(), "gophp_dev_*.phpt")
	if err != nil {
		return "", err
	}
	defer fs.Close()

	return fs.Name(), nil
}

func (m *ExecuteManager) runTestCase(src string, tc *tests.TestCase) (result *tests.Result, log string) {
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

func (m *ExecuteManager) runTestCaseAndReturn(src string, tc *tests.TestCase) *TestResult {
	result, log := m.runTestCase(src, tc)
	sections := tc.Sections()
	statusText := string(result.MainType())

	var status TestResultStatus
	switch statusText {
	case "PASS":
		status = TestResultPass
	case "BORK":
		status = TestResultBork
	case "SKIP":
		status = TestResultSkip
	default:
		status = TestResultFail
	}

	return &TestResult{
		// case
		Code:     sections["FILE"],
		Expected: sections["EXPECT"] + sections["EXPECTF"] + sections["EXPECTREGEX"],

		// result
		Status:     status,
		StatusText: statusText,
		Output:     result.Output(),
		Info:       result.Info() + "\n" + log,
		UseTime:    result.UseTime().Nanoseconds(),
	}
}
