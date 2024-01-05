package instrumenter

import (
	"fmt"
	"os"
	"path/filepath"
)

type MainInstrumentation struct {
	*DefaultInstrumentation
}

func NewMainInstrumentation(pkgPath string, packageBuildDir string) Instrumenter {
	return &MainInstrumentation{
		DefaultInstrumentation: NewDefaultInstrumentation(pkgPath, packageBuildDir).(*DefaultInstrumentation),
	}
}

func (e *MainInstrumentation) WriteExtraFiles(pkgName string) ([]string, error) {
	links := filepath.Join(e.packageBuildDir, "_evomaster_hooks.go")
	err := os.WriteFile(
		links,
		[]byte(
			fmt.Sprintf(
				`package main

import _ "unsafe"
import "github.com/jcbasso/EvoMaster/client-go/src/instrumentation"

//evomaster:ignore

//go:linkname _evomaster_RegisterTargets _evomaster_RegisterTargets
var _evomaster_RegisterTargets = func(ids []string) {
	instrumentation.RegisterTargets(ids)
}

//go:linkname _evomaster_EnteringStatement _evomaster_EnteringStatement
var _evomaster_EnteringStatement = func(fileName string, line int, statement int) {
	instrumentation.EnteringStatement(fileName, line, statement)
}

//go:linkname _evomaster_CompletedStatement _evomaster_CompletedStatement
var _evomaster_CompletedStatement = func(fileName string, line int, statement int) {
	instrumentation.CompletedStatement(fileName, line, statement)
}

//go:linkname _evomaster_CompletionStatement _evomaster_CompletionStatement
var _evomaster_CompletionStatement = func(fileName string, line int, statement int) {
	instrumentation.CompletionStatement(fileName, line, statement)
}

//go:linkname _evomaster_Not _evomaster_Not
var _evomaster_Not = func(value bool) bool {
	return instrumentation.Not(value)
}

//go:linkname _evomaster_And _evomaster_And
var _evomaster_And = func(left func() bool, right func() bool, fileName string, line int, branchId int) bool {
	return instrumentation.And(left, right, fileName, line, branchId)
}

//go:linkname _evomaster_Or _evomaster_Or
var _evomaster_Or = func(left func() bool, right func() bool, fileName string, line int, branchId int) bool {
	return instrumentation.Or(left, right, fileName, line, branchId)
}

//go:linkname _evomaster_CmpUnordered _evomaster_CmpUnordered
var _evomaster_CmpUnordered = func(left any, op string, right any, fileName string, line int, branchId int) bool {
	return instrumentation.CmpUnordered(left, op, right, fileName, line, branchId)
}

//go:linkname _evomaster_CmpOrdered _evomaster_CmpOrdered
var _evomaster_CmpOrdered = func(left any, op string, right any, fileName string, line int, branchId int) bool {
	return instrumentation.CmpOrdered(left, op, right, fileName, line, branchId)
}

%s
`,
				e.registerTargetsCall(),
			),
		),
		0644,
	)
	if err != nil {
		return nil, err
	}
	return []string{links}, nil
}
