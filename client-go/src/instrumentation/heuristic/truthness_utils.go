package heuristic

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
)

// NormalizeValue scales to a positive double value to the [0,1] range.
// v is a non-negative float
func NormalizeValue(v float64) float64 {
	if v < 0 {
		panic(fmt.Sprintf("Negative value: %v", v))
	}

	//normalization function from old ICST/STVR paper
	normalized := v / (v + 1)

	return normalized
}

func GetEqualityTruthness(left any, right any, tracer Tracer) *Truthness {
	lvalue := reflect.ValueOf(left)
	rvalue := reflect.ValueOf(right)

	if !knownPair(lvalue, rvalue) {
		eq := equalityWithNilHandling(left, right, lvalue, rvalue)

		ofTrue := float64(1)
		ofFalse := H_REACHED
		if !eq {
			ofTrue, ofFalse = ofFalse, ofTrue
		}

		return NewTruthness(ofTrue, ofFalse)
	}

	distance := GetDistanceToEquality(lvalue, rvalue, tracer)
	normalizedDistance := NormalizeValue(distance)

	ofFalse := H_REACHED
	if left != right {
		ofFalse = 1
	}
	return NewTruthness(1-normalizedDistance, ofFalse)
}

// Should handle nil separately since interface{}(nil) != struct{}(nil), though nil should be typed in comparison
func equalityWithNilHandling(left, right any, lvalue, rvalue reflect.Value) bool {
	if isNil(left, lvalue) && isNil(right, rvalue) {
		return true
	} else if isNil(left, lvalue) || isNil(right, rvalue) {
		return false
	} else {
		return left == right
	}
}

// Should handle nil separately since interface{}(nil) != struct{}(nil), though nil should be typed in comparison
func isNil(a any, v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return a == nil
	}
}

func GetLessThanTruthness[T constraints.Ordered](left T, right T) *Truthness {
	lvalue := reflect.ValueOf(left)
	rvalue := reflect.ValueOf(right)

	if !knownPair(lvalue, rvalue) {
		lss := left < right

		ofTrue := float64(1)
		ofFalse := H_REACHED
		if !lss {
			ofTrue, ofFalse = ofFalse, ofTrue
		}

		return NewTruthness(ofTrue, ofFalse)
	}

	distance := GetDistanceToLessThan(lvalue, rvalue)

	ofTrue := float64(1)
	ofFalse := 1 / (1.1 + distance)
	if left >= right {
		ofTrue, ofFalse = ofFalse, ofTrue
	}
	return NewTruthness(ofTrue, ofFalse)
}

func knownPair(lvalue reflect.Value, rvalue reflect.Value) bool {
	return BothInt(lvalue, rvalue) || BothUint(lvalue, rvalue) || BothFloat(lvalue, rvalue) || BothString(lvalue, rvalue) || BothIntables(lvalue, rvalue) || BothFloatables(lvalue, rvalue)
}

func BothInt(lvalue reflect.Value, rvalue reflect.Value) bool {
	return lvalue.CanInt() && rvalue.CanInt()
}

func BothUint(lvalue reflect.Value, rvalue reflect.Value) bool {
	return lvalue.CanUint() && rvalue.CanUint()
}

func BothFloat(lvalue reflect.Value, rvalue reflect.Value) bool {
	return lvalue.CanFloat() && rvalue.CanFloat()
}

func BothString(lvalue reflect.Value, rvalue reflect.Value) bool {
	return lvalue.Kind() == reflect.String && rvalue.Kind() == reflect.String
}

func BothIntables(lvalue reflect.Value, rvalue reflect.Value) bool {
	return (lvalue.CanInt() && rvalue.CanUint()) || (lvalue.CanUint() && rvalue.CanInt())
}

func BothFloatables(lvalue reflect.Value, rvalue reflect.Value) bool {
	return (lvalue.CanFloat() && (rvalue.CanInt() || rvalue.CanUint())) || ((lvalue.CanInt() || lvalue.CanUint()) && rvalue.CanFloat())
}
