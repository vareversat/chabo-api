package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsOK(t *testing.T) {
	s := []string{"toto"}
	value := contains(s, "toto")

	assert.True(t, value)
}

func TestContainsNOK(t *testing.T) {
	s := []string{"toto"}
	value := contains(s, "tutu")

	assert.False(t, value)
}

func TestRemoveIOne(t *testing.T) {
	s := []string{"toto", "tutu"}
	expected := []string{"toto"}
	value := remove(s, "tutu")

	assert.True(t, reflect.DeepEqual(expected, value))

}

func TestRemoveITwo(t *testing.T) {
	s := []string{"tutu"}
	expected := []string{}
	value := remove(s, "tutu")

	assert.True(t, reflect.DeepEqual(expected, value))
}

func TestIndexOfOK(t *testing.T) {
	s := []string{"tutu"}
	expected := 0
	value := indexOf(s, "tutu")

	assert.Equal(t, expected, value)
}

func TestIndexOfNOK(t *testing.T) {
	s := []string{"tutu"}
	expected := -1
	value := indexOf(s, "toto")

	assert.Equal(t, expected, value)
}
