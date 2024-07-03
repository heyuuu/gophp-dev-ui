package runservice

var defaultCompileManager Manager = new(CompileManager)
var defaultExecuteManager Manager = new(ExecuteManager)

func GetManager(mode RunMode) Manager {
	switch mode {
	case RunModeCompile:
		return defaultCompileManager
	case RunModeExecute:
		return defaultExecuteManager
	default:
		return nil
	}
}

type Manager interface {
	// run
	AllResultTypes() []ResultType
	Run(code string) *RunResult
	// test
	DefaultTestRoot() string
	FindTestPaths(root string) ([]string, error)
	FindTestCases(root string, path string) ([]string, error)
}
