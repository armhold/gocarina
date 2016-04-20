package gocarina

import (
	"testing"
	"reflect"
)

func TestCanMakeWordFrom(t *testing.T) {
	var examples = []struct {
		word  string
		chars string
		out   bool
	}{
		{"launch", "launch", true},
		{"lunch", "launch", true},
		{"brunch", "launch", false},
		{"aaa", "aaa", true},
		{"aaaa", "aaa", false},
	}

	for _, tt := range examples {
		expected := tt.out
		actual := CanMakeWordFrom(tt.word, tt.chars)

		if actual != expected {
			t.Errorf("error for %q, %q:  wanted %t, got: %t", tt.word, tt.chars, expected, actual)
		}
	}
}

func TestWordsFrom(t *testing.T) {
	expected := []string {
		"ab", "ae", "ar", "arb", "are", "ba",
		"bar", "bare", "be", "bear", "bra", "brae",
		"ea", "ear", "er", "era", "re", "reb",
	}

	actual := WordsFrom("BEAR")

	if ! reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %q, actual: %q", expected, actual)
	}
}
