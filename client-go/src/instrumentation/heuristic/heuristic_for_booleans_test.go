package heuristic_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/heuristic"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate"
	"github.com/stretchr/testify/assert"
	"go/token"
	"math"
	"testing"
)

func TestHeuristicForBooleans_EvaluateOr_TrivialCases(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		leftFunc               func() bool
		rightFunc              func() bool
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:      "true || true",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return true },
			expected:  true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:      "true || false",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return false },
			expected:  true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:      "false || true",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return false },
			expected:  true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:      "false || false",
			leftFunc:  func() bool { return false },
			rightFunc: func() bool { return false },
			expected:  false,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := givenHeuristic.EvaluateOr(c.leftFunc, c.rightFunc, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}

func TestHeuristicForBooleans_EvaluateOr_Chained(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId1 := 1
	givenBranchId2 := 2
	givenBranchId3 := 3
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	leftFunc := func() bool {
		return givenHeuristic.EvaluateOr(func() bool { return false }, func() bool { return false }, givenFileName, givenLine, givenBranchId1, givenTracer)
	}
	rightFunc := func() bool {
		return givenHeuristic.EvaluateOr(func() bool { return false }, func() bool { return true }, givenFileName, givenLine, givenBranchId2, givenTracer)
	}

	// When
	res := givenHeuristic.EvaluateOr(leftFunc, rightFunc, givenFileName, givenLine, givenBranchId3, givenTracer)

	// Then
	expected := true
	assert.Equal(t, expected, res)
	resThenBranch3 := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId3, true)
	resElseBranch3 := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId3, false)
	resThenValue3 := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch3].Value
	resElseValue3 := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch3].Value
	assert.Equal(t, float64(1), resThenValue3)

	expectedElseValue := 0.75
	diff := math.Abs(resElseValue3 - expectedElseValue)
	errMargin := 0.05
	assert.LessOrEqual(t, diff, errMargin)

	givenTracer.Reset()
	givenObjectiveRecorder.Reset(true)
}

func TestHeuristicForBooleans_EvaluateAnd_TrivialCases(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		leftFunc               func() bool
		rightFunc              func() bool
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:      "true && true",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return true },
			expected:  true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:      "true && false",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return false },
			expected:  false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:      "false && true",
			leftFunc:  func() bool { return true },
			rightFunc: func() bool { return false },
			expected:  false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:      "false && false",
			leftFunc:  func() bool { return false },
			rightFunc: func() bool { return false },
			expected:  false,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := givenHeuristic.EvaluateAnd(c.leftFunc, c.rightFunc, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}

func TestHeuristicForBooleans_EvaluateAnd_Chained(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId1 := 1
	givenBranchId2 := 2
	givenBranchId3 := 3
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	leftFunc := func() bool {
		return givenHeuristic.EvaluateAnd(func() bool { return true }, func() bool { return true }, givenFileName, givenLine, givenBranchId1, givenTracer)
	}
	rightFunc := func() bool {
		return givenHeuristic.EvaluateAnd(func() bool { return true }, func() bool { return true }, givenFileName, givenLine, givenBranchId2, givenTracer)
	}

	// When
	res := givenHeuristic.EvaluateAnd(leftFunc, rightFunc, givenFileName, givenLine, givenBranchId3, givenTracer)

	// Then
	expected := true
	assert.Equal(t, expected, res)
	resThenBranch3 := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId3, true)
	resElseBranch3 := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId3, false)
	resThenValue3 := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch3].Value
	resElseValue3 := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch3].Value
	assert.LessOrEqual(t, resThenValue3, float64(1))

	roughExpected := heuristic.H_REACHED // Should be reached but a bit more since it was doubled reached by other branches
	diff := math.Abs(resElseValue3 - roughExpected)
	errMargin := 0.05
	assert.LessOrEqual(t, diff, errMargin)

	givenTracer.Reset()
	givenObjectiveRecorder.Reset(true)
}

func TestHeuristicForBooleans_EvaluateUnorderedCmp(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		left                   any
		op                     token.Token
		right                  any
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:     "int(1) == int(1)",
			left:     1,
			op:       token.EQL,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:     "int(1) == int(10)",
			left:     1,
			op:       token.EQL,
			right:    10,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.1
				diff := math.Abs(value - expected)
				errMargin := 0.0001
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "int(1) != int(1)",
			left:     1,
			op:       token.NEQ,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "int(1) != int(10)",
			left:     1,
			op:       token.NEQ,
			right:    10,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.1
				diff := math.Abs(value - expected)
				errMargin := 0.0001
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "float64(1) == float64(1)",
			left:     float64(1),
			op:       token.EQL,
			right:    float64(1),
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:     "float64(1) == float64(10)",
			left:     float64(1),
			op:       token.EQL,
			right:    float64(10),
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.1
				diff := math.Abs(value - expected)
				errMargin := 0.0001
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "uint(1) == uint(1)",
			left:     uint(1),
			op:       token.EQL,
			right:    uint(1),
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:     "uint(1) == uint(10)",
			left:     uint(1),
			op:       token.EQL,
			right:    uint(10),
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.1
				diff := math.Abs(value - expected)
				errMargin := 0.0001
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "'string' == 'string",
			left:     "string",
			op:       token.EQL,
			right:    "string",
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.LessOrEqual(t, value, heuristic.H_REACHED)
				assert.Greater(t, value, float64(0))
			},
		},
		{
			name:     "'string' == 'stringo'",
			left:     "string",
			op:       token.EQL,
			right:    "stringo",
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.001
				diff := math.Abs(value - expected)
				errMargin := 0.001
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name: "struct == struct",
			left: struct {
				a int
			}{},
			op: token.EQL,
			right: struct {
				a int
			}{},
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := heuristic.H_REACHED
				diff := math.Abs(value - expected)
				errMargin := 0.001
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name: "struct1 == struct2",
			left: struct {
				a int
			}{},
			op: token.EQL,
			right: struct {
				a int64
			}{},
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := heuristic.H_REACHED
				diff := math.Abs(value - expected)
				errMargin := 0.001
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := givenHeuristic.EvaluateUnorderedCmp(c.left, c.op.String(), c.right, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}

func TestHeuristicForBooleans_EvaluateOrderedCmp_Int(t *testing.T) {
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		left                   int64
		op                     token.Token
		right                  int64
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:     "int(2) > int(1)",
			left:     2,
			op:       token.GTR,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "int(1) > int(1)",
			left:     1,
			op:       token.GTR,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.9
				diff := math.Abs(value - expected)
				errMargin := 0.01
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "int(1) < int(1)",
			left:     1,
			op:       token.LSS,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.9
				diff := math.Abs(value - expected)
				errMargin := 0.01
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "int(1) < int(2)",
			left:     1,
			op:       token.LSS,
			right:    2,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "int(2) >= int(1)",
			left:     2,
			op:       token.GEQ,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "int(1) >= int(2)",
			left:     1,
			op:       token.GEQ,
			right:    2,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "int(1) <= int(2)",
			left:     1,
			op:       token.LEQ,
			right:    2,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "int(2) <= int(1)",
			left:     2,
			op:       token.LEQ,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := heuristic.EvaluateOrderedCmp[int64](givenHeuristic, c.left, c.op.String(), c.right, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}

func TestHeuristicForBooleans_EvaluateOrderedCmp_Float(t *testing.T) {
	// Only testing one of each case since the other ones are tested in Integers test
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		left                   float64
		op                     token.Token
		right                  float64
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:     "float64(2) > float64(1)",
			left:     2,
			op:       token.GTR,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "float64(1) < float64(1)",
			left:     1,
			op:       token.LSS,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.9
				diff := math.Abs(value - expected)
				errMargin := 0.01
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "float64(2) >= float64(1)",
			left:     2,
			op:       token.GEQ,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "float64(1) <= float64(2)",
			left:     1,
			op:       token.LEQ,
			right:    2,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := heuristic.EvaluateOrderedCmp[float64](givenHeuristic, c.left, c.op.String(), c.right, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}

func TestHeuristicForBooleans_EvaluateOrderedCmp_Uint(t *testing.T) {
	// Only testing one of each case since the other ones are tested in Integers test
	// Given
	givenHeuristic := heuristic.NewHeuristicForBooleans()
	givenFileName := "file"
	givenLine := 1
	givenBranchId := 1
	givenTracer := staticstate.NewExecutionTracer()
	defer givenTracer.Reset()
	givenObjectiveRecorder := staticstate.NewObjectiveRecorder()
	defer givenObjectiveRecorder.Reset(true)

	cases := []struct {
		name                   string
		left                   uint64
		op                     token.Token
		right                  uint64
		expected               bool
		expectedThenValidation func(*testing.T, float64)
		expectedElseValidation func(*testing.T, float64)
	}{
		{
			name:     "uint64(2) > uint64(1)",
			left:     2,
			op:       token.GTR,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "uint64(1) < uint64(1)",
			left:     1,
			op:       token.LSS,
			right:    1,
			expected: false,
			expectedThenValidation: func(t *testing.T, value float64) {
				expected := 0.9
				diff := math.Abs(value - expected)
				errMargin := 0.01
				assert.LessOrEqual(t, diff, errMargin)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
		},
		{
			name:     "uint64(2) >= uint64(1)",
			left:     2,
			op:       token.GEQ,
			right:    1,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
		{
			name:     "uint64(1) <= uint64(2)",
			left:     1,
			op:       token.LEQ,
			right:    2,
			expected: true,
			expectedThenValidation: func(t *testing.T, value float64) {
				assert.Equal(t, float64(1), value)
			},
			expectedElseValidation: func(t *testing.T, value float64) {
				expected := 0.5
				diff := math.Abs(value - expected)
				errMargin := 0.05
				assert.LessOrEqual(t, diff, errMargin)
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// When
			res := heuristic.EvaluateOrderedCmp[uint64](givenHeuristic, c.left, c.op.String(), c.right, givenFileName, givenLine, givenBranchId, givenTracer)

			// Then
			assert.Equal(t, c.expected, res)
			resThenBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, true)
			resElseBranch := shared.BranchObjectiveName(givenFileName, givenLine, givenBranchId, false)
			resThenValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resThenBranch].Value
			resElseValue := givenTracer.GetInternalReferenceToObjectiveCoverage()[resElseBranch].Value
			c.expectedThenValidation(t, resThenValue)
			c.expectedElseValidation(t, resElseValue)

			givenTracer.Reset()
			givenObjectiveRecorder.Reset(true)
		})
	}
}
