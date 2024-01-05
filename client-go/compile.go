package main

import (
	"errors"
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/instrumenter"
	"log"
	"path/filepath"
	"strings"
)

type compileFlagSet struct {
	Package string `evoflag:"-p"`
	Output  string `evoflag:"-o"`
}

func (f *compileFlagSet) IsValid() bool {
	return f.Package != "" && f.Output != ""
}

func (f *compileFlagSet) String() string {
	return fmt.Sprintf("-p=%q -o=%q", f.Package, f.Output)
}

func ParseCompileCommand(args []string) (commandExecutionFunc, error) {
	if len(args) == 0 {
		return nil, errors.New("unexpected number of command arguments")
	}
	flags := &compileFlagSet{}
	src.ParseFlags(flags, args[1:])

	return makeCompileCommandExecutionFunc(flags, args), nil
}

func makeCompileCommandExecutionFunc(flags *compileFlagSet, args []string) commandExecutionFunc {
	return func() ([]string, error) {
		if !flags.IsValid() {
			// Skip when the required set of flags is not valid.
			return nil, nil
		}

		pkgPath := flags.Package
		packageBuildDir := filepath.Dir(flags.Output)

		i := selectInstrumenter(pkgPath, packageBuildDir)
		if i.IsIgnored() {
			return nil, nil
		}

		argIndices := parseCompileCommandArgs(args)
		log.Println("[args:", args, "]")
		log.Println("[path:", pkgPath, ", buildDir:", packageBuildDir, "]")
		for src := range argIndices {
			log.Println("[src:", src, "]")
		}
		res, err := instrument(i, args, pkgPath, packageBuildDir)
		if err != nil {
			log.Println(err)
		}
		return res, err
	}
}

func selectInstrumenter(pkgPath string, packageBuildDir string) instrumenter.Instrumenter {
	if pkgPath == "main" {
		return instrumenter.NewMainInstrumentation(pkgPath, packageBuildDir)
	}
	if isFromPackage(pkgPath, globalFlags.PathPrefix) {
		return instrumenter.NewDefaultInstrumentation(pkgPath, packageBuildDir)
	}

	return instrumenter.NewIgnoreInstrumentation()
}

// updateArgs Update the argument list by replacing source files that were instrumented.
func updateArgs(args []string, argIndices map[string]int, written map[string]string) {
	for src, dest := range written {
		argIndex := argIndices[src]
		args[argIndex] = dest
	}
}

// parseCompileCommandArgs Walk the list of arguments and add the go source files and the arg slice
// index to returned map.
func parseCompileCommandArgs(args []string) map[string]int {
	goFiles := make(map[string]int)
	for i, src := range args {
		// Only consider args ending with the Go file extension and assume they
		// are Go files.
		if !strings.HasSuffix(src, ".go") {
			continue
		}
		// Save the position of the source file in the argument list
		// to later change it if it gets instrumented.
		goFiles[src] = i
	}
	return goFiles
}

func instrument(i instrumenter.Instrumenter, args []string, pkgPath, packageBuildDir string) ([]string, error) {
	//log.Println("[instrumenting package:", pkgPath, "]")
	//log.Println("[package build directory:", packageBuildDir, "]")

	// Make the list of Go files to instrument out of the argument list and
	// replace their argument list entry by their instrumented copy.
	argIndices := parseCompileCommandArgs(args)
	for src := range argIndices {
		log.Println("[src:", src, "]")
		if err := i.AddFile(src); err != nil {
			return nil, err
		}
	}

	instrumented, err := i.Instrument()
	if err != nil {
		return nil, err
	}

	log.Println("[instrumented_files_#:", len(instrumented), "]")
	pkgName := ""
	if len(instrumented) > 0 {
		pkgName = instrumented[0].Name.Name
		written, err := i.WriteInstrumentedFiles(packageBuildDir, instrumented)
		if err != nil {
			return nil, err
		}
		log.Println(written)
		// Replace original files in the args by the new ones
		updateArgs(args, argIndices, written)
	}

	extraFiles, err := i.WriteExtraFiles(pkgName)
	if err != nil {
		return nil, err
	}

	args = append(args, extraFiles...)
	return args, nil
}

// isFromPackage returns whether the package path is from the main application that is being tested
func isFromPackage(pkgPath string, pathPrefix string) bool {
	if pathPrefix == "" {
		return false
	}

	return strings.HasPrefix(pkgPath, pathPrefix)
}
