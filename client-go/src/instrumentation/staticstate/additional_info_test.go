package staticstate_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/string_specialization"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/taint_type"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdditionalInfo_AddSpecialization(t *testing.T) {
	// Given
	givenAi := staticstate.NewAdditionalInfo()
	givenInfo := string_specialization.StringSpecializationInfo{
		StringSpecialization: string_specialization.EQUAL,
		Value:                "some value",
		TaintType:            taint_type.FULL_MATCH,
	}
	givenTaintName := taint_type.GetTaintName(42)

	// When
	givenAi.AddSpecialization(givenTaintName, givenInfo)
	givenAi.AddSpecialization(givenTaintName, givenInfo)

	// Then
	assert.Len(t, givenAi.StringSpecializations, 1)
	assert.Equal(t, givenAi.StringSpecializations[givenTaintName][givenInfo], true)
}
