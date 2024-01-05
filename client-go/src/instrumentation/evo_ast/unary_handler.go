package evo_ast

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"go/token"
)

type UnaryHandler struct {
	pkgPath string
}

func NewUnaryHandler(pkgPath string) *UnaryHandler {
	return &UnaryHandler{
		pkgPath: pkgPath,
	}
}

func (b *UnaryHandler) Handle(expr *dst.UnaryExpr, pos token.Position, cursor *dstutil.Cursor, fileName string) bool {
	if expr.Op != token.NOT {
		//only handling negation, for now at least...
		return true
	}

	cursor.Replace(b.notExpr(expr.X)) // TODO: Validate if replacement work as intended

	return true
}

func (b *UnaryHandler) notExpr(expr dst.Expr) *dst.ExprStmt {
	// _evomaster_Not(expr)
	return &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun:  dst.NewIdent("_evomaster_Not"),
			Args: []dst.Expr{expr},
		},
	}
}
