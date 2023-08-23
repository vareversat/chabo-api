package utils

import (
	"reflect"
	"testing"
)

func TestContainsOK(t *testing.T) {
	s := []string{"toto"}
	want := true
	value := contains(s, "toto")
	if want != value {
		t.Fatalf(`contains("toto") = %t, want match for %t`, value, want)
	}
}

func TestContainsNOK(t *testing.T) {
	s := []string{"toto"}
	want := false
	value := contains(s, "tutu")
	if want != value {
		t.Fatalf(`contains("toto") = %t, want match for %t`, value, want)
	}
}

func TestRemoveIOne(t *testing.T) {
	s := []string{"toto", "tutu"}
	want := []string{"toto"}
	value := remove(s, "tutu")
	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`remove("tutu") = %q, want match for %#q`, value, want)
	}

}

func TestRemoveITwo(t *testing.T) {
	s := []string{"tutu"}
	want := []string{}
	value := remove(s, "tutu")
	if !reflect.DeepEqual(want, value) {
		t.Fatalf(`remove("tutu") = %q, want match for %#q`, value, want)
	}
}

func TestIndexOfOK(t *testing.T) {
	s := []string{"tutu"}
	want := 0
	value := indexOf(s, "tutu")
	if want != value {
		t.Fatalf(`indexOf("tutu") = %d, want match for %d`, value, want)
	}
}

func TestIndexOfNOK(t *testing.T) {
	s := []string{"tutu"}
	want := -1
	value := indexOf(s, "toto")
	if want != value {
		t.Fatalf(`indexOf("tutu") = %d, want match for %d`, value, want)
	}
}
