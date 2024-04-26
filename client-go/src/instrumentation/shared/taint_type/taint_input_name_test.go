package taint_type_test

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/taint_type"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInputName_Base(t *testing.T) {
	name := taint_type.GetTaintName(0)
	assert.True(t, taint_type.IsTaintInput(name))
}

func TestInputName_InvalidNames(t *testing.T) {
	assert.False(t, taint_type.IsTaintInput("foo"))
	assert.False(t, taint_type.IsTaintInput(""))
	assert.False(t, taint_type.IsTaintInput("evomaster"))
	assert.False(t, taint_type.IsTaintInput("evomaster_input"))
	assert.False(t, taint_type.IsTaintInput("evomaster__input"))
	assert.False(t, taint_type.IsTaintInput("evomaster_a_input"))

	assert.True(t, taint_type.IsTaintInput("_EM_42_XYZ_"))
}

func TestInputName_InvalidNamePatterns(t *testing.T) {
	assert.False(t, taint_type.IsTaintInput("foo"))
	assert.False(t, taint_type.IsTaintInput(""))
	assert.False(t, taint_type.IsTaintInput(taint_type.PREFIX))
	assert.False(t, taint_type.IsTaintInput(taint_type.PREFIX+taint_type.POSTFIX))
	assert.False(t, taint_type.IsTaintInput(taint_type.PREFIX+"a"+taint_type.POSTFIX))

	assert.True(t, taint_type.IsTaintInput(taint_type.PREFIX+"42"+taint_type.POSTFIX))
}

func TestInputName_Includes(t *testing.T) {
	name := taint_type.GetTaintName(0)
	text := "some prefix " + name + " some postfix"

	assert.False(t, taint_type.IsTaintInput(text))
	assert.True(t, taint_type.IncludesTaintInput(text))
}

func TestInputName_UpperLowerCase(t *testing.T) {
	name := taint_type.GetTaintName(0)

	assert.True(t, taint_type.IsTaintInput(name))
	assert.True(t, taint_type.IncludesTaintInput(name))

	assert.True(t, taint_type.IsTaintInput(strings.ToLower(name)))
	assert.True(t, taint_type.IncludesTaintInput(strings.ToLower(name)))
	assert.True(t, taint_type.IsTaintInput(strings.ToUpper(name)))
	assert.True(t, taint_type.IncludesTaintInput(strings.ToUpper(name)))
}
