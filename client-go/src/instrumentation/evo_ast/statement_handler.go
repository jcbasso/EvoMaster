package evo_ast

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
)

type StatementHandler struct {
	pkgPath          string
	statementCounter int64
	objectives       *[]string
}

func NewStatementHandler(pkgPath string, objectives *[]string) *StatementHandler {
	return &StatementHandler{
		pkgPath:          pkgPath,
		statementCounter: 0,
		objectives:       objectives,
	}
}

func (s *StatementHandler) Handle(stmt dst.Stmt, pos token.Position, cursor *dstutil.Cursor, fileName string, fileSet *token.FileSet, astNode ast.Node) bool {
	switch stmt.(type) {
	// TODO: Should add go routines and defer here?
	case *dst.BlockStmt, *dst.ExprStmt: // no point in instrumenting them. Recall, we still instrument its content anyway
		return true
	}

	if cursor.Index() < 0 { // This means that there is no space to insert before and after
		return true
	}

	s.statementCounter++

	*s.objectives = append(*s.objectives, shared.LineObjectiveName(fileName, pos.Line))
	*s.objectives = append(*s.objectives, shared.StatementObjectiveName(fileName, pos.Line, int(s.statementCounter)))

	switch stmt.(type) {
	// TODO: Should add switch?
	case *dst.CaseClause: // ignore each case clause and switch big clause.
		return true
	case *dst.SwitchStmt, *dst.TypeSwitchStmt:
		cursor.InsertBefore(s.completionExpr(pos.Line, fileName, s.statementCounter))
	case *dst.BranchStmt: // continue, break, goto, fallthrough. TODO: Should call it on fallthrough and goto?
		cursor.InsertBefore(s.completionExpr(pos.Line, fileName, s.statementCounter))
	case *dst.ForStmt, *dst.RangeStmt, *dst.IfStmt:
		cursor.InsertBefore(s.completionExpr(pos.Line, fileName, s.statementCounter))
	case *dst.ReturnStmt:
		// TODO: Should do the same for all Return Statement? It returns a function
		cursor.InsertBefore(s.completionExpr(pos.Line, fileName, s.statementCounter))
	default:
		cursor.InsertBefore(s.enteringExpr(pos.Line, fileName, s.statementCounter))
		cursor.InsertAfter(s.completedExpr(pos.Line, fileName, s.statementCounter))
	}

	return true
}

func (s *StatementHandler) completionExpr(line int, fileName string, statementCounter int64) *dst.ExprStmt {
	// _evomaster_CompletionStatement(pkgPath:fileName, line, statementCounter)
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: dst.NewIdent("_evomaster_CompletionStatement"),
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s:%s\"", s.pkgPath, filepath.Base(fileName)), // file name
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Itoa(line), // line
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(statementCounter, 10), // statement
				},
			},
		},
	}
}

func (s *StatementHandler) enteringExpr(line int, fileName string, statementCounter int64) *dst.ExprStmt {
	// _evomaster_EnteringStatement(pkgPath:fileName, line, statementCounter)
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: dst.NewIdent("_evomaster_EnteringStatement"),
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s:%s\"", s.pkgPath, filepath.Base(fileName)), // file name
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Itoa(line), // line
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(statementCounter, 10), // statement
				},
			},
		},
	}
}

func (s *StatementHandler) completedExpr(line int, fileName string, statementCounter int64) *dst.ExprStmt {
	// _evomaster_CompletedStatement(pkgPath:fileName, line, statementCounter)
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: dst.NewIdent("_evomaster_CompletedStatement"),
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s:%s\"", s.pkgPath, filepath.Base(fileName)), // file name
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Itoa(line), // line
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(statementCounter, 10), // statement
				},
			},
		},
	}
}

// TODO: Should be used? It is uncertain the
func (s *StatementHandler) completingExpr(line int, fileName string, statementCounter int64) *dst.ExprStmt {
	// _evomaster_CompletionStatement(pkgPath:fileName, line, statementCounter)
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: dst.NewIdent("_evomaster_CompletionStatement"),
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s:%s\"", s.pkgPath, filepath.Base(fileName)), // file name
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Itoa(line), // line
				},
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(statementCounter, 10), // statement
				},
			},
		},
	}
}
