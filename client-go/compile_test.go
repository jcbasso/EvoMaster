package main

import (
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src"
	"math/rand"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	randomPackageStr := randUTF8String()
	randomOutputStr := randUTF8String()
	randomStr := randUTF8String()

	tests := []struct {
		name                    string
		args                    []string
		expectedFlags           compileFlagSet
		expectedValidationError bool
	}{
		{
			name:                    "empty args",
			args:                    []string{},
			expectedFlags:           compileFlagSet{},
			expectedValidationError: true,
		},
		{
			name: "package option",
			args: []string{fmt.Sprintf("-p=%s", randomPackageStr)},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
			},
			expectedValidationError: true,
		},
		{
			name: "package option",
			args: []string{"-p", randomPackageStr},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
			},
			expectedValidationError: true,
		},
		{
			name: "output option",
			args: []string{fmt.Sprintf("-o=%s", randomOutputStr)},
			expectedFlags: compileFlagSet{
				Output: randomOutputStr,
			},
			expectedValidationError: true,
		},
		{
			name: "output option",
			args: []string{"-o", randomOutputStr},
			expectedFlags: compileFlagSet{
				Output: randomOutputStr,
			},
			expectedValidationError: true,
		},
		{
			name: "output and package options",
			args: []string{"-p", randomPackageStr, "-o", randomOutputStr},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
				Output:  randomOutputStr,
			},
		},
		{
			name: "output and package options",
			args: []string{fmt.Sprintf("-p=%s", randomPackageStr), fmt.Sprintf("-o=%s", randomOutputStr)},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
				Output:  randomOutputStr,
			},
		},
		{
			name: "output and package options",
			args: []string{"-o", randomOutputStr, fmt.Sprintf("-p=%s", randomPackageStr)},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
				Output:  randomOutputStr,
			},
		},
		{
			name: "output and package options",
			args: []string{fmt.Sprintf("-o=%s", randomOutputStr), "-p", randomPackageStr},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
				Output:  randomOutputStr,
			},
		},
		{
			name: "output and package options and others",
			args: []string{fmt.Sprintf("-o=%s", randomOutputStr), "-p", randomPackageStr, "-a", "-b", fmt.Sprintf("-c=%s", randomStr), "-d", randomStr, "a.go", "b.go", "c.go"},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
				Output:  randomOutputStr,
			},
		},
		{
			name: "empty output option value",
			args: []string{"-p", randomPackageStr, "-o", ""},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
			},
			expectedValidationError: true,
		},
		{
			name: "empty package option value",
			args: []string{"-o", randomOutputStr, "-p", ""},
			expectedFlags: compileFlagSet{
				Output: randomOutputStr,
			},
			expectedValidationError: true,
		},
		{
			name: "empty output option value",
			args: []string{"-p", randomPackageStr, "-o="},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
			},
			expectedValidationError: true,
		},
		{
			name: "empty package option value",
			args: []string{"-o", randomOutputStr, "-p="},
			expectedFlags: compileFlagSet{
				Output: randomOutputStr,
			},
			expectedValidationError: true,
		},
		{
			name: "empty package option value",
			args: []string{"-p", "", "-o", randomOutputStr},
			expectedFlags: compileFlagSet{
				Output: randomOutputStr,
			},
			expectedValidationError: true,
		},
		{
			name: "missing output option value",
			args: []string{"-p", randomPackageStr, "-o"},
			expectedFlags: compileFlagSet{
				Package: randomPackageStr,
			},
			expectedValidationError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var flags compileFlagSet
			src.ParseFlags(&flags, tt.args)
			require.Equal(t, tt.expectedFlags, flags)

			isValid := flags.IsValid()
			if tt.expectedValidationError {
				require.False(t, isValid)
			} else {
				require.True(t, isValid)
			}
		})
	}
}

func randUTF8String(length ...int) string {
	numChars := randStringLength(length...)
	codePoints := make([]rune, numChars)
	for i := 0; i < numChars; i++ {
		// Get a random utf8 character code point which is any value between 0 and
		// 0x10FFFF, including non-printable and control characters
		codePoints[i] = rune(rand.Intn(unicode.MaxRune))
	}
	return string(codePoints)
}

func randStringLength(size ...int) (length int) {
	var from, to int
	if len(size) == 1 {
		// String length up to the given max length value
		to = size[0]
	} else if len(size) == 2 {
		// String length between the given boundaries
		from = size[0]
		to = size[1]
	} else {
		// String length up to 1024 characters
		to = rand.Intn(1024)
	}
	upTo := to - from
	if upTo == 0 {
		// Avoid Intn(0) which panics
		return from
	}
	return from + rand.Intn(upTo)

}
