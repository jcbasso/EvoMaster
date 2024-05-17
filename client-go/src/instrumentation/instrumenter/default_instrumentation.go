package instrumenter

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	evo_ast2 "github.com/jcbasso/EvoMaster/client-go/src/instrumentation/evo_ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DefaultInstrumentation struct {
	parsedFiles       map[string]*dst.File
	parsedFileSources map[*dst.File]string
	fset              *token.FileSet
	dstDecorator      *decorator.Decorator
	pkgPath           string
	packageBuildDir   string
	astHandler        *evo_ast2.ASTHandler
	objectives        *[]string
}

func NewDefaultInstrumentation(pkgPath string, packageBuildDir string) Instrumenter {
	objectives := []string{}
	fset := token.NewFileSet() // The token fileset is required to later create the package node.
	dstDecorator := decorator.NewDecorator(fset)
	return &DefaultInstrumentation{
		pkgPath:         unvendorPackagePath(pkgPath),
		packageBuildDir: packageBuildDir,
		astHandler: evo_ast2.NewASTHandler(
			dstDecorator,
			evo_ast2.NewFileHandler(pkgPath, &objectives),
			evo_ast2.NewStatementHandler(pkgPath, &objectives),
			evo_ast2.NewUnaryHandler(pkgPath),
			evo_ast2.NewBinaryHandler(pkgPath, &objectives),
		),
		objectives:   &objectives,
		fset:         fset,
		dstDecorator: dstDecorator,
	}
}

func (e *DefaultInstrumentation) IsIgnored() bool {
	return false
}

func (e *DefaultInstrumentation) AddFile(src string) error {
	// Check if the instrumentation should be skipped for this filename
	if isFileNameIgnored(src) {
		log.Println("skipping instrumentation of file", src)
		return nil
	}

	astFile, err := parser.ParseFile(e.fset, src, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	file, err := e.dstDecorator.DecorateFile(astFile)
	if err != nil {
		return err
	}

	// Check if there is a file-level ignore directive
	ignoreDirective := hasIgnoreDirective(file)
	if ignoreDirective {
		log.Printf("file `%s` skipped due to ignore directive\n", src)
		return nil
	}
	//log.Println("creating parsed files")
	if e.parsedFiles == nil {
		e.parsedFiles = make(map[string]*dst.File)
		e.parsedFileSources = make(map[*dst.File]string)
	}
	e.parsedFiles[src] = file
	e.parsedFileSources[file] = src
	log.Println("[#e.parsedFiles:", len(e.parsedFiles), "]")
	return nil
}

func (e *DefaultInstrumentation) Instrument() ([]*dst.File, error) {
	if len(e.parsedFiles) == 0 {
		log.Println("nothing to instrument")
		return nil, nil
	}

	var instrumentedFiles []*dst.File
	var instrumentedFile *dst.File
	for src, file := range e.parsedFiles {
		//wd, err := os.Getwd()
		//if err != nil {
		//	panic(err)
		//}
		//
		//relSrc, err := filepath.Rel(wd, src) // Using relative src to make it smaller
		//if err != nil {
		//	relSrc = src
		//}

		e.astHandler.SetFileSource(src)
		e.astHandler.SetFileSet(e.fset)

		instrumentedFile = dstutil.Apply(
			file,
			e.astHandler.Pre,
			e.astHandler.Post,
		).(*dst.File) // Safe since file is of type *dst.File

		instrumentedFiles = append(instrumentedFiles, instrumentedFile)
	}

	return instrumentedFiles, nil
}

func (e *DefaultInstrumentation) WriteInstrumentedFiles(packageBuildDir string, instrumentedFiles []*dst.File) (srcast map[string]string, err error) {
	srcast = make(map[string]string, len(instrumentedFiles))
	log.Println("writing files #", len(instrumentedFiles))
	for _, file := range instrumentedFiles {
		src := e.parsedFileSources[file]
		filename := filepath.Base(src)
		dest := filepath.Join(packageBuildDir, filename)
		log.Printf("[src: %s, dst: %s]\n", src, dest)
		output, err := os.Create(dest)
		if err != nil {
			return nil, err
		}
		defer output.Close()

		// Add a go line directive in order to map it to its original source file.
		// Note that otherwise it uses the build directory but it is trimmed by the
		// compiler - so you end up with filenames without any leading path (eg.
		// myfile.go) leading to broken debuggers or stack traces.
		output.WriteString(fmt.Sprintf("//line %s:1\n", src)) // TODO: Validate how
		if err := writeFile(file, output); err != nil {
			return nil, err
		}

		srcast[src] = dest
	}
	return srcast, nil
}

func writeFile(file *dst.File, w io.Writer) error {
	fset, af, err := decorator.RestoreFile(file)
	if err != nil {
		return err
	}
	return printer.Fprint(w, fset, af)
}

func (e *DefaultInstrumentation) WriteExtraFiles(pkgName string) ([]string, error) {
	if pkgName == "" {
		return []string{}, nil
	}
	file := filepath.Join(e.packageBuildDir, "_evomaster_hook.go")

	err := os.WriteFile(
		file,
		[]byte(
			fmt.Sprintf(
				`package %s

import _ "unsafe"

//evomaster:ignore

//go:linkname _evomaster_RegisterTargets _evomaster_RegisterTargets
var _evomaster_RegisterTargets func(ids []string)

//go:linkname _evomaster_EnteringStatement _evomaster_EnteringStatement
var _evomaster_EnteringStatement func(fileName string, line int, statement int)

//go:linkname _evomaster_CompletedStatement _evomaster_CompletedStatement
var _evomaster_CompletedStatement func(fileName string, line int, statement int)

//go:linkname _evomaster_CompletionStatement _evomaster_CompletionStatement
var _evomaster_CompletionStatement func(fileName string, line int, statement int)

//go:linkname _evomaster_Not _evomaster_Not
var _evomaster_Not func(value bool) bool

//go:linkname _evomaster_And _evomaster_And
var _evomaster_And func(left func() bool, right func() bool, fileName string, line int, branchId int) bool

//go:linkname _evomaster_Or _evomaster_Or
var _evomaster_Or func(left func() bool, right func() bool, fileName string, line int, branchId int) bool

//go:linkname _evomaster_CmpUnordered _evomaster_CmpUnordered
var _evomaster_CmpUnordered func(left any, op string, right any, fileName string, line int, branchId int) bool

//go:linkname _evomaster_CmpOrdered _evomaster_CmpOrdered
var _evomaster_CmpOrdered func(left any, op string, right any, fileName string, line int, branchId int) bool

%s
`,
				pkgName,
				e.registerTargetsCall(),
			),
		),
		0644,
	)

	if err != nil {
		return nil, err
	}
	return []string{file}, nil
}

func (e *DefaultInstrumentation) registerTargetsCall() string {
	objectivesStrings := make([]string, len(*e.objectives))
	for i, objective := range *e.objectives {
		objectivesStrings[i] = fmt.Sprintf("\t\t\"%s\",\n", objective)
	}

	return fmt.Sprintf(`
func init() {
	_evomaster_RegisterTargets([]string{
%s
	})
}
`,
		strings.Join(objectivesStrings, ""),
	)
}

// Given the Go vendoring conventions, return the package prefix of the vendored
// package. For example, given `my-app/vendor/github.com/sqreen/go-agent`,
// the function should return `my-app/vendor/`
func unvendorPackagePath(pkg string) (unvendored string) {
	vendorDir := "/vendor/"
	i := strings.Index(pkg, vendorDir)
	if i == -1 {
		return pkg
	}
	return pkg[i+len(vendorDir):]
}

// hasIgnoreDirective Return true if the node has a evomaster:ignore directive comment. Explanatory
// text can be added after it (eg. `//evomaster:ignore because...`)
func hasIgnoreDirective(node dst.Node) bool {
	for _, comment := range node.Decorations().Start.All() {
		if strings.HasPrefix(comment, ignoreDirective) {
			return true
		}
	}
	return false
}

func isFileNameIgnored(file string) bool {
	filename := filepath.Base(file)
	// Don't instrument cgo files
	if strings.Contains(filename, "cgo") {
		return true
	}
	// Don't instrument the go module table file.
	if filename == "_gomod_.go" {
		return true
	}
	return false
}

//// usesImport responds if the file is using the specific import or not
//func usesImport(f *ast.File, name string) (used bool) {
//	ast.Walk(visitFn(func(n ast.Node) {
//		sel, ok := n.(*ast.SelectorExpr)
//		if ok && isTopName(sel.X, name) {
//			used = true
//		}
//	}), f)
//
//	return
//}
//
//type visitFn func(node ast.Node)
//
//func (fn visitFn) Visit(node ast.Node) ast.Visitor {
//	fn(node)
//	return fn
//}
//
//// isTopName returns true if n is a top-level unresolved identifier with the given name.
//func isTopName(n ast.Expr, name string) bool {
//	id, ok := n.(*ast.Ident)
//
//	return ok && id.Name == name && id.Obj == nil
//}
