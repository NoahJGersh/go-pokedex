package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "what    IS up",
			expected: []string{"what", "is", "up"},
		},
		{
			input:    "9 8 765 4.3",
			expected: []string{"9", "8", "765", "4.3"},
		},
		{
			input:    "new\n  line",
			expected: []string{"new", "line"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf(
				"cleanInput failed: len(actual) = %v, len(expected) = %v",
				len(actual),
				len(c.expected),
			)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf(
					"cleanInput failed: actual = %s, expected = %s",
					word,
					expectedWord,
				)
			}
		}
	}
}
