package heuristic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDistanceToEqualityString(t *testing.T) {
	assert.Equal(
		t,
		GetDistanceToEqualityString("eqrt", "sqrt"),
		GetDistanceToEqualityString("sqrt", "eqrt"),
	)
}
