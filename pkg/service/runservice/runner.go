package runservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/heyuuu/gophp/compile/ast"
	"github.com/heyuuu/gophp/compile/ir"
	"github.com/heyuuu/gophp/compile/parser"
	"github.com/heyuuu/gophp/compile/render"
	"github.com/heyuuu/gophp/compile/transformer"
	"github.com/heyuuu/gophp/kits/vardumper"
	_ "github.com/heyuuu/gophp/php/boot"
	"github.com/heyuuu/gophp/php/perr"
	"github.com/heyuuu/gophp/sapi/cli"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RunCode(code string, resultTypes []ResultType) *RunResult {
	r := newRunner(code, resultTypes)
	r.run()

	// build result
	result := new(RunResult)
	result.Result = r.resultItems
	if r.err != nil {
		result.Error = r.err.Error()
	}

	return result
}

type runner struct {
	code          string
	resultTypes   []ResultType
	resultTypeSet map[ResultType]bool

	// result
	resultItems []RunResultItem
	err         error

	// runtime cache
	isTestCase bool
	srcCode    string
	expected   string
	parseRaw   string
	astFile    *ast.File
	irFile     *ir.File
}

func newRunner(code string, resultTypes []ResultType) *runner {
	return &runner{
		code:        code,
		resultTypes: resultTypes,
	}
}

func (r *runner) run() {
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
	isTestCaseMode := strings.HasPrefix(r.code, "!!!")
	if isTestCaseMode {
		testCaseFile := strings.TrimSpace(r.code[3:])
		src, expected, err := loadTestCase(testCaseFile)
		if err != nil {
			r.err = err
			return
		}

		r.isTestCase = true
		r.srcCode = src
		r.expected = expected
	} else {
		r.srcCode = r.code
	}

	// 遍历获取结果
	r.resultItems = make([]RunResultItem, 0, len(r.resultTypes))
	for _, resultType := range r.resultTypes {
		lang, content, ok := r.calcResult(resultType)
		if ok {
			r.resultItems = append(r.resultItems, RunResultItem{
				Type:     resultType,
				Language: lang,
				Content:  content,
			})
		}
	}
}

func (r *runner) checkError(prefix string, err error) {
	if err != nil {
		panic(fmt.Errorf(prefix+" %w", err))
	}
}

func (r *runner) calcResult(typ ResultType) (lang string, content string, ok bool) {
	switch typ {
	case RunTypeSrc:
		return "php", r.srcCode, true
	case RunTypeExpected:
		return "go", r.expected, r.isTestCase
	case RunTypeParseRaw:
		parseRaw, _ := r.requireAstFileEx()
		return "json", parseRaw, true
	case RunTypeAst:
		astFile := r.requireAstFile()
		astDump := vardumper.DumpEx(astFile, vardumper.Config{
			ShowLineNum:        true,
			ShowAnonymousField: true,
		})
		return "", astDump, true
	case RunTypeAstPrint:
		astFile := r.requireAstFile()
		astPrint, err := ast.PrintFile(astFile)
		r.checkError(`ast print fail: `, err)
		return "php", astPrint, true
	case RunTypeIr:
		irFile := r.requireIrFile()
		irDump := vardumper.Dump(irFile)
		return "", irDump, true
	case RunTypeIrPrint:
		irFile := r.requireIrFile()
		irPrint, err := ir.PrintFile(irFile)
		r.checkError("ir print fail: ", err)
		return "php", irPrint, true
	case RunTypeIrRender:
		irFile := r.requireIrFile()
		irRender, err := render.RenderFile(irFile)
		r.checkError("ir render fail: ", err)
		return "go", irRender, true
	case RunTypeExec:
		output := r.executeCode(r.srcCode)
		return "", output, true
	case RunTypeExecRaw:
		output := r.executeCodeRaw(r.srcCode)
		return "", output, true
	default:
		panic(fmt.Errorf("unknown run type: %v", typ))
	}
}

func (r *runner) requireAstFile() *ast.File {
	_, astFile := r.requireAstFileEx()
	return astFile
}

func (r *runner) requireAstFileEx() (string, *ast.File) {
	if r.astFile == nil {
		parseRaw, astFile, err := parser.ParseCodeVerbose(r.srcCode)
		r.checkError(`ast parse fail: `, err)
		r.parseRaw = parseRaw
		r.astFile = astFile
	}
	return r.parseRaw, r.astFile
}

func (r *runner) requireIrFile() *ir.File {
	if r.irFile == nil {
		astFile := r.requireAstFile()
		irFile, err := transformer.TransformFile(astFile)
		r.checkError("ir transformer fail: ", err)
		r.irFile = irFile
	}
	return r.irFile
}

// ResultType: RunTypeExec
func (r *runner) executeCode(code string) string {
	if strings.HasPrefix(code, "<?php\n") {
		code = code[6:]
	} else {
		code = "?>" + code
	}

	cmd := cli.Command("php",
		// ini
		"-d", "error_reporting="+strconv.Itoa(int(perr.E_ALL)),
		// code
		"-r", code,
	)

	var buf strings.Builder
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.RunSafe()
	if err != nil {
		buf.WriteString(">>> Execute failed: " + err.Error())
	}
	return buf.String()
}

// ResultType: RunTypeExecRaw
func (r *runner) executeCodeRaw(code string) string {
	if strings.HasPrefix(code, "<?php\n") {
		code = code[6:]
	} else {
		code = "?>" + code
	}

	output, err := r.runCommand(5*time.Second, "php74", "-r", code)
	if err != nil {
		return output + "\n" + err.Error()
	}

	// todo: 暂时屏蔽行号信息的差异，待完善后移除
	output = regexp.MustCompile(`in Command line code on line \d+`).ReplaceAllString(output, "in Command line code on line 0")

	return output
}

func (r *runner) runCommand(timeout time.Duration, name string, args ...string) (string, error) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	log.Printf("run command: %s\n", cmd.String())
	if output, err := cmd.Output(); err == nil {
		return string(output), nil
	} else if ctx.Err() != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return string(output), fmt.Errorf("run timeout: %w", err)
	} else {
		return string(output), fmt.Errorf("run fail: %w", err)
	}
}
