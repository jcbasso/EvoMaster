package instrumenter

import (
	"github.com/dave/dst"
)

type IgnoreInstrumentation struct {
}

func NewIgnoreInstrumentation() Instrumenter {
	return IgnoreInstrumentation{}
}

func (i IgnoreInstrumentation) IsIgnored() bool {
	return true
}

func (i IgnoreInstrumentation) AddFile(src string) error {
	return nil
}

func (i IgnoreInstrumentation) Instrument() ([]*dst.File, error) {
	return nil, nil
}

func (i IgnoreInstrumentation) WriteInstrumentedFiles(packageBuildDir string, instrumented []*dst.File) (map[string]string, error) {
	return nil, nil
}

func (i IgnoreInstrumentation) WriteExtraFiles(pkgName string) ([]string, error) {
	return nil, nil
}
