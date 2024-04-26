package string_specialization

import (
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/taint_type"
)

type StringSpecializationInfo struct {
	StringSpecialization StringSpecialization

	/**
	 * Value A possible Value to provide context to the specialization.
	 * For example, if the specialization is a CONSTANT, then the "Value" here would
	 * the content of the constant
	 */
	Value string

	TaintType taint_type.TaintType
}

func NewStringSpecializationInfo(
	stringSpecialization StringSpecialization,
	value string,
	taintType taint_type.TaintType, // TODO: Default FULL_MATCH?
) StringSpecializationInfo {
	if taintType <= taint_type.NONE || taintType >= taint_type.End { // NONE shouldn't be possible
		panic(fmt.Sprintf("invalid taint type: %s", taintType.String()))
	}

	return StringSpecializationInfo{
		StringSpecialization: stringSpecialization,
		Value:                value,
		TaintType:            taintType,
	}
}
