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
	declHandler      *DeclHandler

	fileSet      *token.FileSet
	fileSource   string
	dstDecorator *decorator.Decorator
}

func NewASTHandler(decorator *decorator.Decorator, fileHandler *FileHandler, statementHandler *StatementHandler, unaryHandler *UnaryHandler, binaryHandler *BinaryHandler, declHandler *DeclHandler) *ASTHandler {
	return &ASTHandler{
		fileHandler:      fileHandler,
		statementHandler: statementHandler,
		unaryHandler:     unaryHandler,
		binaryHandler:    binaryHandler,
		declHandler:      declHandler,
		dstDecorator:     decorator,
	}
}

func (a *ASTHandler) Pre(c *dstutil.Cursor) bool {
	n := c.Node()
	if n == nil {
		return true // Nothing to do
	}

	switch x := n.(type) {
	case dst.Decl:
		return a.declHandler.Handle(x)
	}

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
		return a.fileHandler.Handle(x, a.fileSource)
	case dst.Stmt:
		return a.statementHandler.Handle(x, pos, c, a.fileSource, a.fileSet, astNode)
	case *dst.BinaryExpr:
		return a.binaryHandler.Handle(x, pos, c, a.fileSource)
	case *dst.UnaryExpr:
		return a.unaryHandler.Handle(x, pos, c, a.fileSource)
	}

	return true
}

func (a *ASTHandler) SetFileSource(src string) {
	a.fileSource = src
}

func (a *ASTHandler) SetFileSet(fset *token.FileSet) {
	a.fileSet = fset
}
