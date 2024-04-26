package taint_type

import (
	"fmt"
	"regexp"
)

/*
   WARNING:
   the naming here has to be kept in sync in ALL implementations of this class,
   including Java, JS and C#
*/

const PREFIX = "_EM_"
const POSTFIX = "_XYZ_"

// TODO: Validate why it has \\Q and \\E in Java but not JS  Pattern.compile("\\Q"+PREFIX+"\\E\\d+\\Q"+POSTFIX+"\\E");
const REGEX = PREFIX + "\\d+" + POSTFIX

func partialMatch() *regexp.Regexp {
	r, err := regexp.Compile("(?i)" + REGEX)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile regular expression %s: %s", REGEX, err.Error()))
	}

	return r
}

// var fullMatch = new RegExp("^" + TaintInputName.regex + "$", "i")
func fullMatch() *regexp.Regexp {
	r, err := regexp.Compile("^" + "(?i)" + REGEX + "$")
	if err != nil {
		panic(fmt.Sprintf("Failed to compile regular expression %s: %s", REGEX, err.Error()))
	}

	return r
}

// IsTaintInput Check if a given string value is a tainted value
func IsTaintInput(value string) bool {
	if value == "" {
		return false
	}

	return fullMatch().MatchString(value)
}

func IncludesTaintInput(value string) bool {
	if value == "" {
		return false
	}

	return partialMatch().MatchString(value)
}

// GetTaintName Create a tainted value, with the input id being part of it
func GetTaintName(id int) string { // TODO: Should be int? Should validate negatives, looks odd to do so?
	/*
	   Note: this is quite simple, we simply add a unique prefix
	   and postfix, in lowercase.
	   But we would not be able to check if the part of the id was
	   modified.
	*/
	return fmt.Sprintf("%s%d%s", PREFIX, id, POSTFIX)
}
