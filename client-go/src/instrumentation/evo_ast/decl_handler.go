package evo_ast

import (
	"github.com/dave/dst"
)

type DeclHandler struct {
	pkgPath string
}

func NewDeclHandler(pkgPath string) *DeclHandler {
	return &DeclHandler{
		pkgPath: pkgPath,
	}
}

func (s *DeclHandler) Handle(decl dst.Decl) bool {
	switch d := decl.(type) {
	case *dst.FuncDecl:
		if d.Name.Name == "init" {
			// Skip the entire subtree of the init function
			return false
		}
	}

	return true
}
