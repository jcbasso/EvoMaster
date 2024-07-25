package heuristic

import (
	"math"
	"reflect"
)

const (
	// MAX_CHAR_DISTANCE 2^16=65536, max distance for a char
	MAX_CHAR_DISTANCE = 65_536

	// TODO: Validate if adding collection distance calculation
	//  H_REACHED_BUT_NULL  = 0.05
	//  H_NOT_NULL          = 0.1
	//  H_REACHED_BUT_EMPTY = H_REACHED_BUT_NULL
	//  H_NOT_EMPTY         = H_NOT_NULL

)

func GetDistanceToEqualityInt64(a int64, b int64) float64 {
	return GetDistanceToEqualityFloat64(float64(a), float64(b))
}

func GetDistanceToEqualityUint64(a uint64, b uint64) float64 {
	return GetDistanceToEqualityFloat64(float64(a), float64(b))
}

func GetDistanceToEqualityFloat64(a float64, b float64) float64 {
	if math.IsInf(a, 0) || math.IsInf(b, 0) || math.IsNaN(a) || math.IsNaN(b) {
		return math.MaxFloat64
	}

	var distance float64 = math.Abs(a - b)
	if distance < 0 {
		// overflow has occurred
		return math.MaxFloat64
	} else {
		return distance
	}
}

func GetDistanceToEqualityString(a string, b string) float64 {
	return GetLeftAlignmentDistance(a, b)
}

func GetLeftAlignmentDistance(a string, b string) float64 {
	diff := math.Abs(float64(len(a) - len(b)))
	dist := diff * MAX_CHAR_DISTANCE

	min := len(a)
	if len(a) > len(b) {
		min = len(b)
	}

	if min == 0 { // One of the strings is empty
		return dist
	}

	for i := 0; i < min; i++ {
		dist += math.Abs(float64(rune(a[i]) - rune(b[i])))
	}

	return dist
}

// GetDistanceToEquality compute distance between left and right ==
func GetDistanceToEquality(lvalue reflect.Value, rvalue reflect.Value, tracer Tracer) float64 {
	if BothInt(lvalue, rvalue) {
		return GetDistanceToEqualityInt64(lvalue.Int(), rvalue.Int())
	}

	if BothUint(lvalue, rvalue) {
		return GetDistanceToEqualityUint64(lvalue.Uint(), rvalue.Uint())
	}

	if BothFloat(lvalue, rvalue) {
		return GetDistanceToEqualityFloat64(lvalue.Float(), rvalue.Float())
	}

	if BothString(lvalue, rvalue) {
		tracer.HandleTaintForStringEquals(lvalue.String(), rvalue.String(), false)
		return GetDistanceToEqualityString(lvalue.String(), rvalue.String())
	}

	return math.MaxFloat64
}

//	TODO: Validate if adding collection distance calculation
///**
// * Return a h=[0,1] heuristics from a scaled distance, taking into account a starting base
// * @param base
// * @param distance
// * @return
// */
//func HeuristicFromScaledDistanceWithBase(base: number, distance: number): number{
//
//	if (base < 0 || base >=1)
//		throw Error("Invalid base: " + base);
//	if (distance < 0)
//		throw Error("Negative distance: " + distance)
//	if (!isFinite(distance) || distance == Number.MAX_VALUE)
//		return base;
//
//	return base + ((1-base)/(distance + 1));
//}

func GetDistanceToLessThanInt64(a int64, b int64) float64 {
	return GetDistanceToLessThanFloat64(float64(a), float64(b))
}

func GetDistanceToLessThanUint64(a uint64, b uint64) float64 {
	return GetDistanceToLessThanFloat64(float64(a), float64(b))
}

func GetDistanceToLessThanFloat64(a float64, b float64) float64 {
	return GetDistanceToEqualityFloat64(a, b)
}

func GetDistanceToLessThanString(a string, b string) float64 {
	distance := float64(MAX_CHAR_DISTANCE)

	for i := 0; i < len(a) && i < len(b); i++ {
		/*
		   What determines the order is the first char they have different,
		   starting from left to right
		*/
		if a[i] == b[i] {
			continue
		}

		distance = math.Abs(float64(a[i] - b[i]))
		break
	}

	return distance
}

// GetDistanceToLessThan compute distance between left and right <
func GetDistanceToLessThan(lvalue reflect.Value, rvalue reflect.Value) float64 {
	if BothInt(lvalue, rvalue) {
		return GetDistanceToLessThanInt64(lvalue.Int(), rvalue.Int())
	}

	if BothUint(lvalue, rvalue) {
		return GetDistanceToLessThanUint64(lvalue.Uint(), rvalue.Uint())
	}

	if BothFloat(lvalue, rvalue) {
		return GetDistanceToLessThanFloat64(lvalue.Float(), rvalue.Float())
	}

	if BothString(lvalue, rvalue) {
		return GetDistanceToLessThanString(lvalue.String(), rvalue.String())
	}

	// Both should be int
	if BothIntables(lvalue, rvalue) {
		if lvalue.CanInt() {
			return GetDistanceToLessThanInt64(lvalue.Int(), int64(rvalue.Uint()))
		} else {
			return GetDistanceToLessThanInt64(int64(lvalue.Uint()), rvalue.Int())
		}
	}

	// Both should be floats
	if BothFloatables(lvalue, rvalue) {
		if lvalue.CanFloat() {
			if rvalue.CanInt() {
				return GetDistanceToLessThanFloat64(lvalue.Float(), float64(rvalue.Int()))
			} else {
				return GetDistanceToLessThanFloat64(lvalue.Float(), float64(rvalue.Uint()))
			}
		} else { // rvalue.CanFloat()
			if lvalue.CanInt() {
				return GetDistanceToLessThanFloat64(float64(lvalue.Int()), rvalue.Float())
			} else {
				return GetDistanceToLessThanFloat64(float64(lvalue.Uint()), rvalue.Float())
			}
		}
	}

	return math.MaxFloat64
}
