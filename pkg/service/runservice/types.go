package runservice

type RunMode string

const (
	RunModeCompile RunMode = "compile"
	RunModeExecute RunMode = "execute"
)

type ResultType string

const (
	RunTypeAst      ResultType = "ast"
	RunTypeAstPrint ResultType = "ast-print"
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
