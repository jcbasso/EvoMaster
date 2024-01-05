package instrumenter

import (
	"github.com/dave/dst"
)

type Instrumenter interface {
	IsIgnored() bool
	AddFile(src string) error
	Instrument() ([]*dst.File, error)
	WriteInstrumentedFiles(packageBuildDir string, instrumented []*dst.File) (srcdst map[string]string, err error)
	WriteExtraFiles(pkgName string) ([]string, error)
}
