package instrumentation_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate"
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func Test_CmpOrdered_Success(t *testing.T) {
	// Given
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	defer staticstate.NewExecutionTracer().Reset()
	defer staticstate.NewObjectiveRecorder().Reset(true)

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
		{
			name:     "int vs uint64",
			value1:   int(0),
			value2:   uint(0),
			op:       token.GEQ,
			expected: true,
		},
		{
			name:     "int vs float",
			value1:   int(0),
			value2:   float64(0),
			op:       token.GEQ,
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// When
			f := func() {
				res := instrumentation.CmpOrdered(c.value1, c.op.String(), c.value2, givenFileName, givenLine, givenBranchId)

				// Then
				assert.Equal(t, c.expected, res)
			}

			// Then
			assert.NotPanics(t, f)
		})
	}
}

func Test_CmpOrdered_Panic_on_difference(t *testing.T) {
	// Given
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenValue1 := int(1)
	givenValue2 := "1"
	givenOp := token.LEQ
	defer staticstate.NewExecutionTracer().Reset()
	defer staticstate.NewObjectiveRecorder().Reset(true)

	// When
	f := func() {
		instrumentation.CmpOrdered(givenValue1, givenOp.String(), givenValue2, givenFileName, givenLine, givenBranchId)
	}

	//Then
	assert.Panics(t, f)
}

func Test_CmpUnordered_Success(t *testing.T) {
	// Given
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	defer staticstate.NewExecutionTracer().Reset()
	defer staticstate.NewObjectiveRecorder().Reset(true)

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
			value2:   int32(1),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "int8 vs int64",
			value1:   int8(1),
			value2:   int64(1),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "uint vs uint32",
			value1:   uint(1),
			value2:   uint32(1),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "uint8 vs uint64",
			value1:   uint8(1),
			value2:   uint64(1),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "float64 vs float32",
			value1:   float64(1),
			value2:   float32(1),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "string",
			value1:   "string",
			value2:   "string",
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "int vs uint64",
			value1:   int(600000),
			value2:   uint(600000),
			op:       token.EQL,
			expected: true,
		},
		{
			name:     "int vs float: correctly",
			value1:   int(0),
			value2:   float64(0.00000000000000006),
			op:       token.EQL,
			expected: false,
		},
		//{
		//	name:     "int vs float: panics since 1-0.00000000000000005 looks to calculate as 1 and Truthness can't be 1 and 1",
		//	value1:   int(0),
		//	value2:   float64(0.00000000000000005),
		//	op:       token.EQL,
		//	expected: true,
		//},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// When
			f := func() {
				res := instrumentation.CmpUnordered(c.value1, c.op.String(), c.value2, givenFileName, givenLine, givenBranchId)
				assert.Equal(t, c.expected, res)
			}

			// Then
			assert.NotPanics(t, f)
		})
	}
}
