package runservice

import (
	"fmt"
	"github.com/heyuuu/gophp/compile/ast"
	"github.com/heyuuu/gophp/compile/ir"
	"github.com/heyuuu/gophp/compile/parser"
	"github.com/heyuuu/gophp/compile/render"
	"github.com/heyuuu/gophp/compile/transformer"
	"github.com/heyuuu/gophp/kits/oskit"
	"github.com/heyuuu/gophp/kits/vardumper"
	"strings"
)

type compileRunner struct {
	resultItems []RunResultItem
	err         error
}

func (r *compileRunner) reset() {
	r.resultItems = nil
	r.err = nil
}

func (r *compileRunner) run(code string) *RunResult {
	// reset
	r.reset()

	// run
	r.realRunCode(code)

	// build result
	result := new(RunResult)
	result.Result = r.resultItems
	if r.err != nil {
		result.Error = r.err.Error()
	}

	return result
}

func (r *compileRunner) realRunCode(code string) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok && err != nil {
				r.err = err
			} else {
				r.err = fmt.Errorf("unknown run panic: %v", e)
			}
		}
	}()

	// load test case
	isTestCaseMode := strings.HasPrefix(code, "!!!")
	if isTestCaseMode {
		testCaseFile := strings.TrimSpace(code[3:])
		src, expected, err := r.loadTestCase(testCaseFile)
		if err != nil {
			r.err = err
			return
		}
		r.addResult(RunTypeSrc, "php", src)
		r.addResult(RunTypeExpected, "go", expected)

		code = src
	}

	r.runCode(code)
}

func (r *compileRunner) loadTestCase(srcFile string) (src string, expected string, err error) {
	if !strings.HasSuffix(srcFile, ".php") {
		err = fmt.Errorf("test case file must be .php")
		return
	}

	src, err = oskit.ReadFileAsString(srcFile)
	if err != nil {
		err = fmt.Errorf("load src file failed: %w", err)
		return
	}

	expectedFile := strings.TrimSuffix(srcFile, ".php") + ".go"
	expected, err = oskit.ReadFileAsString(expectedFile)
	if err != nil {
		err = fmt.Errorf("load expected file failed: %w", err)
		return
	}

	return
}

func (r *compileRunner) checkError(prefix string, err error) {
	if err != nil {
		panic(fmt.Errorf(prefix+" %w", err))
	}
}

func (r *compileRunner) addResult(typ ResultType, lang string, content string) {
	r.resultItems = append(r.resultItems, RunResultItem{
		Type:     typ,
		Language: lang,
		Content:  content,
	})
}

func (r *compileRunner) runCode(code string) {
	// parse
	parseRaw, astFile, err := parser.ParseCodeVerbose(code)
	r.addResult(RunTypeParseRaw, "json", parseRaw)
	r.checkError(`ast parse fail: `, err)

	astDump := vardumper.DumpEx(astFile, vardumper.Config{
		ShowLineNum:        true,
		ShowAnonymousField: true,
	})
	r.addResult(RunTypeAst, "", astDump)

	astPrint, err := ast.PrintFile(astFile)
	r.checkError(`ast print fail: `, err)
	r.addResult(RunTypeAstPrint, "php", astPrint)

	// transformer
	irFile, err := transformer.TransformFile(astFile)
	r.checkError("ir transformer fail: ", err)

	irDump := vardumper.Dump(irFile)
	r.addResult(RunTypeIr, "", irDump)

	irPrint, err := ir.PrintFile(irFile)
	r.checkError("ir print fail: ", err)
	r.addResult(RunTypeIrPrint, "php", irPrint)

	// render
	irRender, err := render.RenderFile(irFile)
	r.checkError("ir render fail: ", err)
	r.addResult(RunTypeIrRender, "go", irRender)
}
