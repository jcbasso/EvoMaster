package main

import (
	"errors"
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var globalFlags src.InstrumentationToolFlagSet

func main() {
	log.SetFlags(0)
	//log.SetPrefix("evomaster: ")
	log.SetOutput(os.Stderr)

	args := os.Args[1:]
	cmd, cmdArgPos, err := parseCommand(&globalFlags, args)

	if err != nil || globalFlags.Help {
		printUsage()
		os.Exit(1)
	}

	// Hide instrumentation tool arguments
	if cmdArgPos != -1 {
		args = args[cmdArgPos:]
	}

	var logs strings.Builder
	//if !globalFlags.Verbose {
	//	// Save the logs to show them in case of instrumentation error
	//	log.SetOutput(&logs)
	//}

	if cmd != nil {
		// The command is implemented
		newArgs, err := cmd()
		if err != nil {
			//log.Println(err)
			//if !globalFlags.Verbose {
			fmt.Fprintln(os.Stderr, &logs)
			//}
			os.Exit(1)
		}
		if newArgs != nil {
			// Args are replaced
			args = newArgs
		}
	}

	err = forwardCommand(args)
	var exitErr *exec.ExitError
	if err != nil {
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatalln(err)
		}
	}
	os.Exit(0)
}

// forwardCommand runs the given command's argument list and exits the process
// with the exit code that was returned.
func forwardCommand(args []string) error {
	path := args[0]
	args = args[1:]
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//quotedArgs := fmt.Sprintf("%+q", args)
	//log.Printf("forwarding command `%s %s`", path, quotedArgs[1:len(quotedArgs)-1])
	return cmd.Run()
}

func printUsage() {
	const usageFormat = `Usage: go {build,install,get,test} -a -toolexec '%s [-path-prefix <path>]' PACKAGES...

Instrumentation tool for Go v%s. It instruments Go source code at
compilation time by adding hooks.

Options:
        -h
                Print this usage message.
		-path-prefix
				Prefix of the path to instrument. Is used to infer which sub-packages are related to your package and 
				instrument them. If there are none subdirectories can omit but probably you want to have it set.

To see the instrumented code, use the go option -work in order to keep the
build directory. It will contain every instrumented Go source file.
`
	_, _ = fmt.Fprintf(os.Stderr, usageFormat, os.Args[0], "1.21")
	os.Exit(2)
}

type parseCommandFunc func([]string) (commandExecutionFunc, error)
type commandExecutionFunc func() (newArgs []string, err error)

var commandParserMap = map[string]parseCommandFunc{
	"compile": ParseCompileCommand,
}

// getCommand returns the command and arguments. The command is expectedFlags to be
// the first argument.
func parseCommand(instrToolFlagSet *src.InstrumentationToolFlagSet, args []string) (commandExecutionFunc, int, error) {
	//log.Println("[cmds: ", args, "]")
	cmdIdPos := src.ParseFlagsUntilFirstNonOptionArg(instrToolFlagSet, args)
	//globalFlags.Full = true
	if cmdIdPos == -1 {
		return nil, cmdIdPos, errors.New("unexpected arguments")
	}
	cmdId := args[cmdIdPos]
	args = args[cmdIdPos:]
	cmdId, err := parseCommandID(cmdId)
	if err != nil {
		return nil, cmdIdPos, err
	}

	if commandParser, exists := commandParserMap[cmdId]; exists {
		cmd, err := commandParser(args)
		return cmd, cmdIdPos, err
	} else {
		return nil, cmdIdPos, nil
	}
}

func parseCommandID(cmd string) (string, error) {
	// It mustn't be empty
	if cmd == "" {
		return "", errors.New("unexpected empty command name")
	}

	// Take the base of the absolute path of the go tool
	cmd = filepath.Base(cmd)
	// Remove the file extension if any
	if ext := filepath.Ext(cmd); ext != "" {
		cmd = strings.TrimSuffix(cmd, ext)
	}
	return cmd, nil
}
