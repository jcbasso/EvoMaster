package shared

import (
	"fmt"
)

// FILE Prefix identifier for file coverage objectives.
// A file is "covered" if at least one of its lines is executed.
//
// Note: this is different from Java where we rather look at CLASS
const FILE = "File"

// LINE Prefix identifier for line coverage objectives
const LINE = "Line"

// STATEMENT Prefix identifier for statement coverage objectives
const STATEMENT = "Statement"

// BRANCH Prefix identifier for branch coverage objectives
const BRANCH = "Branch"

// TRUE_BRANCH Tag used in a branch id to specify it is for the "true"/then branch
const TRUE_BRANCH = "_trueBranch"

// FALSE_BRANCH Tag used in a branch id to specify it is for the "false"/else branch
const FALSE_BRANCH = "_falseBranch"

func FileObjectiveName(fileID string) string {
	return fmt.Sprintf("%s_%s", FILE, fileID)
}

func GetFileIdFromObjectiveName(target string) string {
	prefix := fmt.Sprintf("%s_", FILE)
	return target[len(prefix):]
}

func LineObjectiveName(fileID string, line int) string {
	return fmt.Sprintf("%s_%s_%s", LINE, fileID, padNumber(line))
}

func StatementObjectiveName(fileID string, line int, index int) string {
	return fmt.Sprintf("%s_%s_%s_%v", STATEMENT, fileID, padNumber(line), index)
}

func BranchObjectiveName(fileID string, line int, branchID int, thenBranch bool) string {
	branch := FALSE_BRANCH
	if thenBranch {
		branch = TRUE_BRANCH
	}
	return fmt.Sprintf("%s_at_%s_at_line_%s_position_%v%s", BRANCH, fileID, padNumber(line), branchID, branch)
}

func padNumber(value int) string {
	if value < 0 {
		panic("negative number to pad")
	}

	return fmt.Sprintf("%05d", value)
}
