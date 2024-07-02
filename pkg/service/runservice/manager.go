package runservice

var defaultExecuteManager Manager = new(ExecuteManager)

func GetManager(mode RunMode) Manager {
	switch mode {
	case RunModeExecute:
		return defaultExecuteManager
	default:
		return nil
	}
}

type Manager interface {
	AllResultTypes() []ResultType
	Run(code string) *RunResult
	FindTestPaths(root string) []string
	FindTestCases(path string) []string
}
