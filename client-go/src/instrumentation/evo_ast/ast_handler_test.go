package evo_ast

import (
	"bytes"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/stretchr/testify/assert"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"
)

func Test_ASTHandler_TypeSwitch(t *testing.T) {
	// Given
	pkgPath := "some-pkg"
	objectives := []string{}
	fset := token.NewFileSet() // The token fileset is required to later create the package node.
	dstDecorator := decorator.NewDecorator(fset)
	astHandler := NewASTHandler(
		dstDecorator,
		NewFileHandler(pkgPath, &objectives),
		NewStatementHandler(pkgPath, &objectives),
		NewUnaryHandler(pkgPath),
		NewBinaryHandler(pkgPath, &objectives),
		NewDeclHandler(pkgPath),
	)
	astHandler.SetFileSet(fset)
	astHandler.SetFileSource("file-source")

	cases := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name: "switch",
			file: `package main
				
				func SwitchTest() bool {
					switch someVariable {
					case someValue:
						return false
					default:
						return true
					}
				}`,
			expected: `package main
				
				func SwitchTest() bool {
					_evomaster_CompletionStatement("some-pkg:file-source", 4, %id%)
					switch someVariable {
					case someValue:
						_evomaster_CompletionStatement("some-pkg:file-source", 6, %id%)
						return false
					default:
						_evomaster_CompletionStatement("some-pkg:file-source", 8, %id%)
						return true
					}
				}`,
		},
		{
			name: "type switch",
			file: `package main
				
				func TypeSwitchTest() bool {
					switch _ := variable.(type) {
					case someType:
						return false
					default:
						return true
					}
				}`,
			expected: `package main
				
				func TypeSwitchTest() bool {
					_evomaster_CompletionStatement("some-pkg:file-source", 4, %id%)
					switch _ := variable.(type) {
					case someType:
						_evomaster_CompletionStatement("some-pkg:file-source", 6, %id%)
						return false
					default:
						_evomaster_CompletionStatement("some-pkg:file-source", 8, %id%)
						return true
					}
				}`,
		},
		{
			name: "init ignore",
			file: `package main
				
				func init() {
					return true
				}

				func otherFunc() {
					return true
				}`,
			expected: `package main
				
				func init() {
					return true
				}

				func otherFunc() {
					_evomaster_CompletionStatement("some-pkg:file-source", 8, %id%)
					return true
				}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			astNode, err := parser.ParseFile(fset, "file", c.file, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}

			dstNode, err := dstDecorator.DecorateFile(astNode)
			if err != nil {
				t.Fatal(err)
			}

			// When
			res := dstutil.Apply(
				dstNode,
				astHandler.Pre,
				astHandler.Post,
			).(*dst.File)

			// Then
			s, err := fileToString(res)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, trim(c.expected), trim(s))
		})
	}
}

func fileToString(file *dst.File) (string, error) {
	fset, astFile, err := decorator.RestoreFile(file)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = format.Node(buf, fset, astFile)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func trim(s string) string {
	tmp := strings.Replace(s, "\t", "", -1)
	tmp = strings.TrimSuffix(tmp, "\n")
	return replaceInjectedParams(tmp)
}

func replaceInjectedParams(code string) string {
	var res string
	var regex *regexp.Regexp
	// _evomaster_EnteringStatement
	regex = regexp.MustCompile(`_evomaster_EnteringStatement\((.*,.*, ?)(.+)\)`)
	res = regex.ReplaceAllString(code, `_evomaster_EnteringStatement($1%id%)`)
	// _evomaster_CompletedStatement
	regex = regexp.MustCompile(`_evomaster_CompletedStatement\((.*,.*, ?)(.+)\)`)
	res = regex.ReplaceAllString(code, `_evomaster_CompletedStatement($1%id%)`)
	// _evomaster_CompletionStatement
	regex = regexp.MustCompile(`_evomaster_CompletionStatement\((.*,.*, ?)(.+)\)`)
	res = regex.ReplaceAllString(code, `_evomaster_CompletionStatement($1%id%)`)

	return res
}
