package evo_ast

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"go/token"
	"path/filepath"
	"strconv"
)

type BinaryHandler struct {
	pkgPath       string
	branchCounter int64
	objectives    *[]string
}

func NewBinaryHandler(pkgPath string, objectives *[]string) *BinaryHandler {
	return &BinaryHandler{
		pkgPath:    pkgPath,
		objectives: objectives,
	}
}

func (b *BinaryHandler) Handle(expr *dst.BinaryExpr, pos token.Position, cursor *dstutil.Cursor, fileName string) bool {
	// valid token
	validToken := true
	switch expr.Op {
	case token.LAND: // &&
		cursor.Replace(b.callBinaryExpr("_evomaster_And", expr.X, expr.Y, fileName, pos.Line, b.branchCounter))
	case token.LOR: // ||
		cursor.Replace(b.callBinaryExpr("_evomaster_Or", expr.X, expr.Y, fileName, pos.Line, b.branchCounter))
	case
		token.EQL, // ==
		token.NEQ: // !=
		cursor.Replace(b.callCmpExpr("_evomaster_CmpUnordered", expr.X, expr.Op, expr.Y, fileName, pos.Line, b.branchCounter))
	case
		token.LSS, // <
		token.LEQ, // <=
		token.GTR, // >
		token.GEQ: // >=
		cursor.Replace(b.callCmpExpr("_evomaster_CmpOrdered", expr.X, expr.Op, expr.Y, fileName, pos.Line, b.branchCounter))
	default:
		validToken = false
	}

	if !validToken {
		return true // nothing to do
	}

	b.branchCounter++

	*b.objectives = append(*b.objectives, shared.BranchObjectiveName(fileName, pos.Line, int(b.branchCounter), true))
	*b.objectives = append(*b.objectives, shared.BranchObjectiveName(fileName, pos.Line, int(b.branchCounter), false))

	return true
}

func (b *BinaryHandler) callBinaryExpr(funcName string, left dst.Expr, right dst.Expr, fileName string, line int, branchCounter int64) dst.Expr {
	// funcName(func() bool { return left }, func() bool { return right }, fileName, line, branchCounter)
	return &dst.CallExpr{
		Fun: dst.NewIdent(funcName),
		Args: []dst.Expr{
			b.toLambda(left),
			b.toLambda(right),
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s:%s\"", b.pkgPath, filepath.Base(fileName)), // file name
			},
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Itoa(line), // line
			},
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: strconv.FormatInt(branchCounter, 10), // branch
			},
		},
	}
}

func (b *BinaryHandler) callCmpExpr(funcName string, left dst.Expr, op token.Token, right dst.Expr, fileName string, line int, branchCounter int64) dst.Expr {
	// funcName(left, op, right, fileName, line, branchCounter)
	return &dst.CallExpr{
		Fun: dst.NewIdent(funcName),
		Args: []dst.Expr{
			left,
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s\"", op.String()), // file name
			},
			right,
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s:%s\"", b.pkgPath, filepath.Base(fileName)), // file name
			},
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Itoa(line), // line
			},
			&dst.BasicLit{
				Kind:  token.STRING,
				Value: strconv.FormatInt(branchCounter, 10), // branch
			},
		},
	}
}

func (b *BinaryHandler) toLambda(expr dst.Expr) *dst.FuncLit {
	return &dst.FuncLit{
		Type: &dst.FuncType{
			Results: &dst.FieldList{
				List: []*dst.Field{{
					Type: &dst.Ident{
						Name: "bool",
					},
				}},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{expr},
				},
			},
		},
	}
}
