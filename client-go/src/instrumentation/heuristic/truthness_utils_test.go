package heuristic_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/heuristic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTruthnessUtils_GetEqualityTruthness(t *testing.T) {
	// Given
	cases := []struct {
		name     string
		left     any
		right    any
		expected *heuristic.Truthness
	}{
		{
			name:     "0 == 0",
			left:     0,
			right:    0,
			expected: heuristic.NewTruthness(1, 0.01),
		},
		{
			name:     "0.0 == 0.0",
			left:     0.0,
			right:    0.0,
			expected: heuristic.NewTruthness(1, 0.01),
		},
		{
			name:     "int32(0) == int32(0)",
			left:     int32(0),
			right:    int32(0),
			expected: heuristic.NewTruthness(1, 0.01),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := heuristic.GetEqualityTruthness(c.left, c.right)

			// Then
			assert.Equal(t, c.expected, res)
		})
	}
}
