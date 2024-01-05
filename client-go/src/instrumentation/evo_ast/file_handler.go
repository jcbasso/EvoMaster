package evo_ast

import (
	"github.com/dave/dst"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
)

type FileHandler struct {
	pkgPath    string
	objectives *[]string
}

func NewFileHandler(pkgPath string, objectives *[]string) *FileHandler {
	return &FileHandler{
		pkgPath:    pkgPath,
		objectives: objectives,
	}
}

func (f *FileHandler) Handle(file *dst.File, fileName string) bool {
	*f.objectives = append(*f.objectives, shared.FileObjectiveName(fileName))

	return true
}
