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

	RunTypeParseRaw ResultType = "parse_raw"
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
