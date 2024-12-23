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

type Runner struct {
	compileMode bool
	resultItems []RunResultItem
	err         error
}

func NewRunner(compileMode bool) *Runner {
	return &Runner{
		compileMode: compileMode,
	}
}

func (r *Runner) reset() {
	r.resultItems = nil
	r.err = nil
}

func (r *Runner) Run(code string) *RunResult {
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

func (r *Runner) realRunCode(code string) {
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
		src, expected, err := loadTestCase(testCaseFile)
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

func (r *Runner) checkError(prefix string, err error) {
	if err != nil {
		panic(fmt.Errorf(prefix+" %w", err))
	}
}

func (r *Runner) addResult(typ ResultType, lang string, content string) {
	r.resultItems = append(r.resultItems, RunResultItem{
		Type:     typ,
		Language: lang,
		Content:  content,
	})
}

func (r *Runner) runCode(code string) {
	if r.compileMode {
		r.runCodeCompile(code)
	} else {
		r.runCodeExecute(code)
	}
}

func (r *Runner) runCodeCompile(code string) {
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

func (r *Runner) runCodeExecute(code string) {
	// parse-ast
	astNodes, err := parser.ParseCode(code)
	r.checkError("ast parse fail: %w", err)

	astDump := vardumper.Dump(astNodes)
	r.addResult(RunTypeAst, "", astDump)

	astPrint, err := ast.PrintFile(astNodes)
	r.checkError("ast print fail: %w", err)
	r.addResult(RunTypeAstPrint, "", astPrint)

	// run code
	output := r.executeCode(code)
	r.addResult(RunTypeExec, "", output)

	// raw run code
	rawOutput := r.executeCodeRaw(code)
	r.addResult(RunTypeExecRaw, "", rawOutput)
}

func (r *Runner) executeCode(code string) string {
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

func (r *Runner) executeCodeRaw(code string) string {
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

func (r *Runner) runCommand(timeout time.Duration, name string, args ...string) (string, error) {
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
