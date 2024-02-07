package slices

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestUnique(t *testing.T) {
	var testcases = []struct {
		input []string

		expected []string
	}{
		{
			input:    nil,
			expected: []string{},
		},
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"a", "c", "b", "a", "c", "a", "d", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range testcases {
		got := Unique(tt.input)
		if slices.Compare(got, tt.expected) != 0 {
			t.Errorf("got %v, want %v", got, tt.expected)
		}
	}
}
