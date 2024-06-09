package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

type RunType = string

const (
	RunTypeAst      RunType = "AST"
	RunTypeAstPrint RunType = "AST-print"
	RunTypeExec     RunType = "Execute"
	RunTypeExecRaw  RunType = "Execute-Raw"
)

var allRunTypes = []RunType{
	RunTypeAst,
	RunTypeAstPrint,
	RunTypeExec,
	RunTypeExecRaw,
}

// api: /run/config
func RunConfigHandler(c *gin.Context) any {
	return gin.H{
		"types": allRunTypes,
	}
}

// api: /run/code
func RunCodeHandler(c *gin.Context) any {
	// 获取请求参数
	var p runCodeParam
	if err := c.ShouldBind(&p); err != nil {
		return err
	}

	// 执行代码，并返回结果
	result, err := runCode(p.Code)
	if err != nil {
		return err
	}

	return gin.H{
		"result": result,
	}
}

type runCodeParam struct {
	Code string `form:"code" binding:"required"`
}

type RunResultItem struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func runCode(code string) (result []RunResultItem, err error) {
	// parse-ast
	astNodes, err := parser.ParseCode(code)
	if err != nil {
		return nil, fmt.Errorf("ast parse fail: %w", err)
	}

	astDump := vardumper.Sprint(astNodes)
	result = append(result, RunResultItem{Type: RunTypeAst, Content: astDump})

	astPrint, err := ast.PrintFile(astNodes)
	if err != nil {
		return nil, fmt.Errorf("ast print fail: %w", err)
	}
	result = append(result, RunResultItem{Type: RunTypeAstPrint, Content: astPrint})

	// run code
	output := executeCode(code)
	result = append(result, RunResultItem{Type: RunTypeExec, Content: output})

	// raw run code
	rawOutput := executeCodeRaw(code)
	result = append(result, RunResultItem{Type: RunTypeExecRaw, Content: rawOutput})

	return
}

func executeCode(code string) string {
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

func executeCodeRaw(code string) string {
	if strings.HasPrefix(code, "<?php\n") {
		code = code[6:]
	} else {
		code = "?>" + code
	}

	output, err := runCommand(5*time.Second, "php74", "-r", code)
	if err != nil {
		return output + "\n" + err.Error()
	}

	// todo: 暂时屏蔽行号信息的差异，待完善后移除
	output = regexp.MustCompile(`in Command line code on line \d+`).ReplaceAllString(output, "in Command line code on line 0")

	return output
}

func runCommand(timeout time.Duration, name string, args ...string) (string, error) {
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
