package src

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindFirstNonOptionArg(t *testing.T) {
	tests := []struct {
		args            []string
		expectedPos     int
		expectedFlagSet InstrumentationToolFlagSet
	}{
		{
			args:        nil,
			expectedPos: -1,
		},
		{
			args:        []string{},
			expectedPos: -1,
		},
		{
			args:        []string{"-", "-", "-"},
			expectedPos: -1,
		},
		{
			args:        []string{"-", "-", "-"},
			expectedPos: -1,
		},
		{
			args:        []string{"-", "c", "-"},
			expectedPos: 1,
		},
		{
			args:        []string{"-", "-", "c"},
			expectedPos: 2,
		},
		{
			args:        []string{"c", "-", "-"},
			expectedPos: 0,
		},
		{
			args:        []string{"-v", "-full", "cmd", "-a", "b", "c"},
			expectedPos: 2,
			expectedFlagSet: InstrumentationToolFlagSet{
				Verbose: true,
				Full:    true,
			},
		},
		{
			args:        []string{"-v", "what", "-ever", "-is=this", "-full", "cmd", "-a", "b", "c"},
			expectedPos: 1,
			expectedFlagSet: InstrumentationToolFlagSet{
				Verbose: true,
			},
		},
		{
			args:        []string{"-v", "-what", "-ever", "-is=this", "-full", "cmd", "-a", "b", "c"},
			expectedPos: 5,
			expectedFlagSet: InstrumentationToolFlagSet{
				Verbose: true,
				Full:    true,
			},
		},

		{
			args:        []string{"-v", "/usr/lib/go-1.13/pkg/tool/linux_amd64/compile", "-V=full"},
			expectedPos: 1,
			expectedFlagSet: InstrumentationToolFlagSet{
				Verbose: true,
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(t *testing.T) {
			var flagSet InstrumentationToolFlagSet
			got := ParseFlagsUntilFirstNonOptionArg(&flagSet, tc.args)
			require.Equal(t, tc.expectedPos, got)
			require.Equal(t, tc.expectedFlagSet, flagSet)
		})
	}
}
