package runservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/heyuuu/gophp/compile/ast"
	"github.com/heyuuu/gophp/compile/parser"
	"github.com/heyuuu/gophp/kits/vardumper"
	_ "github.com/heyuuu/gophp/php/boot"
	"github.com/heyuuu/gophp/php/perr"
	"github.com/heyuuu/gophp/sapi"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type executeRunner struct {
	result *RunResult
}

func (r *executeRunner) run(code string) (result *RunResult) {
	result = new(RunResult)
	r.result = result

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok && err != nil {
				result.Error = err.Error()
			} else {
				result.Error = fmt.Sprintf("unknown run panic: %v", e)
			}
		}
	}()

	r.runCode(code)

	return
}

func (r *executeRunner) checkError(prefix string, err error) {
	if err != nil {
		panic(fmt.Errorf(prefix+" %w", err))
	}
}

func (r *executeRunner) addResult(typ ResultType, lang string, content string) {
	r.result.Result = append(r.result.Result, RunResultItem{
		Type:     typ,
		Language: lang,
		Content:  content,
	})
}

func (r *executeRunner) runCode(code string) {
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

func (r *executeRunner) executeCode(code string) string {
	if strings.HasPrefix(code, "<?php\n") {
		code = code[6:]
	} else {
		code = "?>" + code
	}

	cmd := sapi.Command(
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

func (r *executeRunner) executeCodeRaw(code string) string {
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

func (r *executeRunner) runCommand(timeout time.Duration, name string, args ...string) (string, error) {
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
