package instrumentation_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/execution_tracer"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/objective_recorder"
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func Test_CmpOrdered_Success(t *testing.T) {
	// Given
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	defer execution_tracer.New().Reset()
	defer objective_recorder.New().Reset(true)

	cases := []struct {
		name     string
		value1   any
		value2   any
		op       token.Token
		expected bool
	}{
		{
			name:     "int vs int32",
			value1:   int(1),
			value2:   int32(2),
			op:       token.LSS,
			expected: true,
		},
		{
			name:     "int8 vs int64",
			value1:   int8(2),
			value2:   int64(1),
			op:       token.LSS,
			expected: false,
		},
		{
			name:     "uint vs uint32",
			value1:   uint(1),
			value2:   uint32(2),
			op:       token.LSS,
			expected: true,
		},
		{
			name:     "uint8 vs uint64",
			value1:   uint8(2),
			value2:   uint64(1),
			op:       token.LSS,
			expected: false,
		},
		{
			name:     "float64 vs float32",
			value1:   float64(1),
			value2:   float32(2),
			op:       token.LSS,
			expected: true,
		},
		{
			name:     "string",
			value1:   "string",
			value2:   "string other",
			op:       token.LSS,
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// When
			res := instrumentation.CmpOrdered(c.value1, c.op.String(), c.value2, givenFileName, givenLine, givenBranchId)

			// Then
			assert.Equal(t, c.expected, res)
		})
	}
}

func Test_CmpOrdered_Panic_on_difference(t *testing.T) {
	// Given
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenValue1 := int(1)
	givenValue2 := uint(1)
	givenOp := token.LEQ
	defer execution_tracer.New().Reset()
	defer objective_recorder.New().Reset(true)

	// When & Then
	assert.Panics(t, func() {
		instrumentation.CmpOrdered(givenValue1, givenOp.String(), givenValue2, givenFileName, givenLine, givenBranchId)
	})
}
