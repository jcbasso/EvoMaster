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

func GetEqualityTruthness(left any, right any) *Truthness {
	lvalue := reflect.ValueOf(left)
	rvalue := reflect.ValueOf(right)

	if !knownPair(lvalue, rvalue) {
		eq := left == right

		ofTrue := float64(1)
		ofFalse := H_REACHED
		if !eq {
			ofTrue, ofFalse = ofFalse, ofTrue
		}

		return NewTruthness(ofTrue, ofFalse)
	}

	distance := GetDistanceToEquality(lvalue, rvalue)
	normalizedDistance := NormalizeValue(distance)

	ofFalse := H_REACHED
	if left != right {
		ofFalse = 1
	}
	return NewTruthness(1-normalizedDistance, ofFalse)
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
	return BothInt(lvalue, rvalue) || BothUint(lvalue, rvalue) || BothFloat(lvalue, rvalue) || BothString(lvalue, rvalue)
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
