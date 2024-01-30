package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartWithUpperCase_OK(t *testing.T) {
	s := "Test"
	value := StartWithUpperCase(s)

	assert.True(t, value)
}

func TestStartWithUpperCase_NOK(t *testing.T) {
	s := "test"
	value := StartWithUpperCase(s)

	assert.False(t, value)
}

func TestEndWithLowerCase_OK(t *testing.T) {
	s := "test"
	value := EndWithLowerCase(s)

	assert.True(t, value)
}

func TestEndWithLowerCase_NOK(t *testing.T) {
	s := "tesT"
	value := EndWithLowerCase(s)

	assert.False(t, value)
}
