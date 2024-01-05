package evo_ast

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/format"
	"testing"
)

func Test_callBinaryExpr(t *testing.T) {
	// Given
	givenFuncName := "myFunction"
	givenLeftExpr := &dst.Ident{Name: "true"}
	givenRightExpr := &dst.Ident{Name: "false"}
	givenFileName := "fileName"
	givenLine := 1
	var givenBranchId int64 = 2
	givenPkgPath := "my_pkg"
	target := NewBinaryHandler(givenPkgPath, nil)

	// When
	res := target.callBinaryExpr(givenFuncName, givenLeftExpr, givenRightExpr, givenFileName, givenLine, givenBranchId)
	resString, err := toString(res)
	assert.Nil(t, err)

	// Then
	expected := fmt.Sprintf(`%s(func() bool {
	return true
}, func() bool {
	return false
}, "%s:%s", %d, %d)`, givenFuncName, givenPkgPath, givenFileName, givenLine, givenBranchId)
	assert.Equal(t, expected, resString)
}

func toString(expr dst.Expr) (string, error) {
	// Convert to DST file since it is the way to convert to AST
	dstFile := &dst.File{
		Name: &dst.Ident{Name: "main"},
		Decls: []dst.Decl{
			&dst.FuncDecl{
				Name: &dst.Ident{Name: "main"},
				Type: &dst.FuncType{},
				Body: &dst.BlockStmt{
					List: []dst.Stmt{
						&dst.ExprStmt{
							X: expr,
						},
					},
				},
			},
		},
	}

	// Convert to AST file
	fset, astFile, err := decorator.RestoreFile(dstFile)
	if err != nil {
		return "", err
	}

	// Obtain only the needed statement
	astExpr := (astFile.Decls[0].(*ast.FuncDecl)).Body.List[0]

	// Parse to string
	buf := new(bytes.Buffer)
	err = format.Node(buf, fset, astExpr)
	if err != nil {
		return "", err
	}
	res := buf.String()

	return res, nil
}
