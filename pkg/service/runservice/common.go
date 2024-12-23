package runservice

import (
	"fmt"
	"github.com/heyuuu/gophp/kits/oskit"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func findTestFiles(root string, path string, checker func(file string) bool) []string {
	var testFiles []string
	_ = eachTestFile(filepath.Join(root, path), checker, func(fullPath string) {
		name, _ := filepath.Rel(root, fullPath)
		testFiles = append(testFiles, name)
	})
	return testFiles
}

func eachTestFile(dir string, checker func(file string) bool, handler func(file string)) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == "" || file.Name()[0] == '.' {
			continue
		}

		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			err = eachTestFile(path, checker, handler)
		} else { // IsFile
			if checker(path) {
				handler(path)
			}
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func findTestPaths(root string, checker func(file string) bool) []string {
	var testPaths []string
	eachTestPath(root, checker, func(dir string) {
		path, _ := filepath.Rel(root, dir)
		testPaths = append(testPaths, path)
	})
	slices.Sort(testPaths)
	return testPaths
}

func eachTestPath(dir string, checker func(file string) bool, handler func(file string)) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	var isTestPath bool
	for _, file := range files {
		if file.Name() == "" || file.Name()[0] == '.' {
			continue
		}

		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if eachTestPath(path, checker, handler) {
				isTestPath = true
			}
		} else { // IsFile
			if !isTestPath && checker(path) {
				isTestPath = true
			}
		}
	}
	if isTestPath {
		handler(dir)
	}

	return isTestPath
}

func loadTestCase(srcFile string) (src string, expected string, err error) {
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
