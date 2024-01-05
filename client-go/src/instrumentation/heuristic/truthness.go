package heuristic

import "fmt"

// 2 values: one for true, and one for false.
// The values are in [0,1].
// One of them is necessarily equal to 1 (which represents the actual result of the expression),
// but not both, ie an expression evaluates to either true or false.
// The non-1 value represents how close the other option would have been from being taken.

type Truthness struct {
	OfTrue  float64
	OfFalse float64
}

func NewTruthness(ofTrue float64, ofFalse float64) *Truthness {
	// TODO: Should panic or should return error?
	if ofTrue < 0 || ofTrue > 1 {
		panic(fmt.Sprintf("invalid value for ofTrue: %f", ofTrue))
	}

	if ofFalse < 0 || ofFalse > 1 {
		panic(fmt.Sprintf("invalid value for ofFalse: %f", ofFalse))
	}

	//NOTE: no longer the case
	//if ofTrue != 1 && ofFalse != 1 {
	//	panic(fmt.Sprintf("at least one value should be equal to 1"))
	//}

	if ofTrue == 1 && ofFalse == 1 {
		panic("values cannot be both equal to 1")
	}

	return &Truthness{
		OfTrue:  ofTrue,
		OfFalse: ofFalse,
	}
}

func (t *Truthness) Invert() *Truthness {
	return NewTruthness(t.OfFalse, t.OfTrue)
}

func (t *Truthness) RescaleFromMin(min float64) *Truthness {
	var ofTrue float64 = 1
	if t.OfTrue != 1 {
		ofTrue = min + (1-min)*t.OfTrue
	}

	var ofFalse float64 = 1
	if t.OfFalse != 1 {
		ofFalse = min + (1-min)*t.OfFalse
	}

	return NewTruthness(ofTrue, ofFalse)
}

func (t *Truthness) IsTrue() bool {
	return t.OfTrue == 1
}

func (t *Truthness) IsFalse() bool {
	return t.OfFalse == 1
}
