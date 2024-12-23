package runservice

type RunMode string

const (
	RunModeCompile RunMode = "compile"
	RunModeExecute RunMode = "execute"
)

type ResultType string

const (
	RunTypeSrc      ResultType = "src"
	RunTypeExpected ResultType = "expected"

	RunTypeParseRaw ResultType = "parse-raw"
	RunTypeAst      ResultType = "ast"
	RunTypeAstPrint ResultType = "ast-print"
	RunTypeIr       ResultType = "ir"
	RunTypeIrPrint  ResultType = "ir-print"
	RunTypeIrRender ResultType = "ir-render"
	RunTypeExec     ResultType = "execute"
	RunTypeExecRaw  ResultType = "execute-raw"
)

type RunResultItem struct {
	Type     ResultType `json:"type"`
	Language string     `json:"language"`
	Content  string     `json:"content"`
}

type RunResult struct {
	Result []RunResultItem `json:"result"`
	Error  string          `json:"error"`
}

type TestResultStatus string

const (
	TestResultPass TestResultStatus = "PASS" // 执行成功
	TestResultFail TestResultStatus = "FAIL" // 执行失败
	TestResultBork TestResultStatus = "BORK" // 测试case不合法
	TestResultSkip TestResultStatus = "SKIP" // 跳过执行
)

type TestResult struct {
	Code       string           `json:"code"`
	Expected   string           `json:"expected"`
	Status     TestResultStatus `json:"status"`
	StatusText string           `json:"status_text"`
	Output     string           `json:"output"`
	Info       string           `json:"info"`
	UseTime    int64            `json:"useTime"`
}
