package evo_ast

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"go/token"
)

type ASTHandler struct {
	fileHandler      *FileHandler
	statementHandler *StatementHandler
	unaryHandler     *UnaryHandler
	binaryHandler    *BinaryHandler

	fileSet      *token.FileSet
	fileSource   string
	dstDecorator *decorator.Decorator
}

func NewASTHandler(decorator *decorator.Decorator, fileHandler *FileHandler, statementHandler *StatementHandler, unaryHandler *UnaryHandler, binaryHandler *BinaryHandler) *ASTHandler {
	return &ASTHandler{
		fileHandler:      fileHandler,
		statementHandler: statementHandler,
		unaryHandler:     unaryHandler,
		binaryHandler:    binaryHandler,
		dstDecorator:     decorator,
	}
}

func (a *ASTHandler) Pre(c *dstutil.Cursor) bool {
	return true
}

func (a *ASTHandler) Post(c *dstutil.Cursor) bool {
	n := c.Node()
	if n == nil {
		return true // Nothing to do
	}

	astNode, ok := a.dstDecorator.Map.Ast.Nodes[n]
	if !ok {
		panic("missing ast node mapping")
	}

	pos := a.fileSet.Position(astNode.Pos())
	switch x := n.(type) {
	case *dst.File:
		a.fileHandler.Handle(x, a.fileSource)
	case dst.Stmt:
		a.statementHandler.Handle(x, pos, c, a.fileSource, a.fileSet, astNode)
	case *dst.BinaryExpr:
		a.binaryHandler.Handle(x, pos, c, a.fileSource)
	case *dst.UnaryExpr:
		a.unaryHandler.Handle(x, pos, c, a.fileSource)
	}

	return true
}

func (a *ASTHandler) SetFileSource(src string) {
	a.fileSource = src
}

func (a *ASTHandler) SetFileSet(fset *token.FileSet) {
	a.fileSet = fset
}
