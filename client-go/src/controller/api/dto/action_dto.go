package dto

type ActionDto struct {
	// The index of this action in the test.
	// Eg, in a test with 10 indices, the index would be
	// between 0 and 9
	Index int `json:"index"`

	// A list (possibly empty) of String values used in the action.
	// This info can be used for different kinds of taint analysis, eg
	// to check how such values are used in the SUT
	InputVariables []string `json:"inputVariables"`
}
