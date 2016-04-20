package gocarina

import (
	"testing"
	"reflect"
	"sort"
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
		"bare", "bear", "brae", "arb", "are", "bar", "bra", "ear", "era", "reb", "ab", "ae", "ar", "ba", "be", "ea", "er", "re",
	}

	actual := WordsFrom("BEAR")

	if ! reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %q, actual: %q", expected, actual)
	}
}

func TestSortByWordLength(t *testing.T) {
	var examples = []struct {
		in  []string
		out []string
	}{
		// bigger words first, but then alphabetical for equal word-lengths
		{[]string{ "door", "dead", "darth", "dear", "apple", "a", "alpha", "zylophone", "beta", "bear"}, []string{ "zylophone", "alpha", "apple", "darth", "bear", "beta", "dead", "dear", "door", "a"}},
	}

	for _, tt := range examples {
		expected := tt.out

		// first make a copy, because sort.Sort() works in-place on the slice
		actual := make([]string, len(tt.in))
		copy(actual, tt.in)
		sort.Sort(ByWordLength(actual))

		if ! reflect.DeepEqual(expected, actual) {
			t.Errorf("error for %q, wanted %q, got: %q", tt.in, expected, actual)
		}
	}
}
