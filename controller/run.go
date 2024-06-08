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

type runCodeParam struct {
	Code string `form:"code" binding:"required"`
}

func ApiRunCode(c *gin.Context) any {
	var err error

	var p runCodeParam
	if err = c.ShouldBind(&p); err != nil {
		return apiError(err)
	}

	result, err := parseCode(p.Code)
	if err != nil {
		return apiError(err)
	}

	return apiSucc(gin.H{
		"result": result,
	})
}

const (
	TypeAst      = "AST"
	TypeAstPrint = "AST-print"
	TypeRun      = "Run"
	TypeRawRun   = "Run-Raw"
)

type ApiTypeResult struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func parseCode(code string) (result []ApiTypeResult, err error) {
	// Ast
	astNodes, err := parser.ParseCode(code)
	if err != nil {
		return nil, fmt.Errorf("ast parse fail: %w", err)
	}

	astDump := vardumper.Sprint(astNodes)
	result = append(result, ApiTypeResult{Type: TypeAst, Content: astDump})

	astPrint, err := ast.PrintFile(astNodes)
	if err != nil {
		return nil, fmt.Errorf("ast print fail: %w", err)
	}
	result = append(result, ApiTypeResult{Type: TypeAstPrint, Content: astPrint})

	// run code
	output := runCode(code)
	result = append(result, ApiTypeResult{Type: TypeRun, Content: output})

	// raw run code
	rawOutput := rawRunCode(code)
	result = append(result, ApiTypeResult{Type: TypeRawRun, Content: rawOutput})

	return
}

func runCode(code string) string {
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

func rawRunCode(code string) string {
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
