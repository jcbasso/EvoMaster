package heuristic

import (
	"fmt"
	"go/token"
	"golang.org/x/exp/constraints"
	"math"
	"sync"
)

const (
	H_REACHED = 0.01 // branch was reached
)

type HeuristicForBooleans struct {
	lastEvaluation   *Truthness
	lastEvaluationMx sync.RWMutex
}

var heuristicForBooleansOnce sync.Once
var heuristicForBooleansInstance *HeuristicForBooleans

func NewHeuristicForBooleans() *HeuristicForBooleans {
	heuristicForBooleansOnce.Do(func() {
		heuristicForBooleansInstance = &HeuristicForBooleans{}
	})

	return heuristicForBooleansInstance
}

func (h *HeuristicForBooleans) ClearLastEvaluation() {
	h.lastEvaluationMx.Lock()
	h.lastEvaluation = nil
	h.lastEvaluationMx.Unlock()
}

func (h *HeuristicForBooleans) HandleNot(value bool) bool {
	h.lastEvaluationMx.Lock()
	if h.lastEvaluation != nil {
		h.lastEvaluation = h.lastEvaluation.Invert()
	}
	h.lastEvaluationMx.Unlock()
	return !value
}

func (h *HeuristicForBooleans) EvaluateAnd(left func() bool, right func() bool, fileName string, line int, branchId int, tracer Tracer) bool {
	// (x && y) == ! (!x || !y)

	h.ClearLastEvaluation()

	var x bool
	var xT *Truthness

	x = left()
	xT = h.calculateTruthness(x)

	var t *Truthness
	var y bool

	if x {

		h.ClearLastEvaluation()

		var yT *Truthness

		y = right()
		yT = h.calculateTruthness(y)

		t = NewTruthness(
			(xT.OfTrue/2)+(yT.OfTrue/2),
			math.Max(xT.OfFalse, yT.OfFalse),
		)
	} else {
		t = NewTruthness(
			xT.OfTrue/2,
			xT.OfFalse,
		)
	}

	tracer.UpdateBranch(fileName, line, branchId, t)

	h.lastEvaluationMx.Lock()
	h.lastEvaluation = t
	h.lastEvaluationMx.Unlock()

	return x && y
}

func (h *HeuristicForBooleans) EvaluateOr(left func() bool, right func() bool, fileName string, line int, branchId int, tracer Tracer) bool {
	h.ClearLastEvaluation()

	var x bool
	var xT *Truthness

	x = left()
	xT = h.calculateTruthness(x)

	var t *Truthness
	var y bool

	if !x {

		h.ClearLastEvaluation()

		var yT *Truthness

		y = right()
		yT = h.calculateTruthness(y)

		t = NewTruthness(
			math.Max(xT.OfTrue, yT.OfTrue),
			(xT.OfFalse/2)+(yT.OfFalse/2),
		)
	} else {
		t = NewTruthness(
			xT.OfTrue,
			xT.OfFalse/2,
		)
	}

	tracer.UpdateBranch(fileName, line, branchId, t)

	h.lastEvaluationMx.Lock()
	h.lastEvaluation = t
	h.lastEvaluationMx.Unlock()

	return x || y
}

func (h *HeuristicForBooleans) calculateTruthness(branch bool) *Truthness {
	base := H_REACHED

	h.lastEvaluationMx.RLock()
	lastEvaluation := h.lastEvaluation
	h.lastEvaluationMx.RUnlock()
	if lastEvaluation == nil {
		ofTrue := float64(1)
		ofFalse := base
		if !branch {
			ofTrue, ofFalse = ofFalse, ofTrue
		}

		return NewTruthness(ofTrue, ofFalse)
	}

	h.lastEvaluationMx.Lock()
	defer h.lastEvaluationMx.Unlock()
	return h.lastEvaluation.RescaleFromMin(base)
}

func (h *HeuristicForBooleans) EvaluateUnorderedCmp(left any, op string, right any, fileName string, line int, branchId int, tracer Tracer) bool {
	/*
	   Make sure we get exactly the same result
	*/
	var res bool
	if op == token.EQL.String() {
		res = left == right
	} else if op == token.NEQ.String() {
		res = left != right
	} else {
		panic(fmt.Sprintf("Invalid op: %s", op))
	}

	t := h.compareUnordered(left, op, right)

	tracer.UpdateBranch(fileName, line, branchId, t)

	h.lastEvaluationMx.Lock()
	h.lastEvaluation = t
	h.lastEvaluationMx.Unlock()

	return res
}

func (h *HeuristicForBooleans) compareUnordered(left any, op string, right any) *Truthness {
	if op == token.EQL.String() {
		return GetEqualityTruthness(left, right)
	} else if op == token.NEQ.String() {
		return h.compareUnordered(left, token.EQL.String(), right).Invert()
	}

	panic(fmt.Sprintf("Invalid op: %s", op))
}

func EvaluateOrderedCmp[T constraints.Ordered](h *HeuristicForBooleans, left T, op string, right T, fileName string, line int, branchId int, tracer Tracer) bool {
	// Using function out of the struct to be able to use Type Generics

	/*
	   Make sure we get exactly the same result
	*/
	var res bool
	if op == token.LSS.String() {
		res = left < right
	} else if op == token.LEQ.String() {
		res = left <= right
	} else if op == token.GTR.String() {
		res = left > right
	} else if op == token.GEQ.String() {
		res = left >= right
	} else {
		panic(fmt.Sprintf("Invalid op: %s", op))
	}

	t := compareOrdered[T](h, left, op, right)

	tracer.UpdateBranch(fileName, line, branchId, t)

	h.lastEvaluationMx.Lock()
	h.lastEvaluation = t
	h.lastEvaluationMx.Unlock()

	return res
}

func compareOrdered[T constraints.Ordered](h *HeuristicForBooleans, left T, op string, right T) *Truthness {
	if op == token.LSS.String() {
		return GetLessThanTruthness(left, right)
		// (l <= r)  same as  (r >= l)  same as  !(r < l)
	} else if op == token.GEQ.String() {
		return compareOrdered(h, left, token.LSS.String(), right).Invert()
	} else if op == token.LEQ.String() {
		return compareOrdered(h, right, token.LSS.String(), left).Invert()
	} else if op == token.GTR.String() {
		return compareOrdered(h, left, token.LEQ.String(), right).Invert()
	}

	panic(fmt.Sprintf("Invalid op: %s", op))
}

type Tracer interface {
	UpdateBranch(fileName string, line int, branch int, truthness *Truthness)
}
