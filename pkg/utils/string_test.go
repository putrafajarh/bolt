package utils_test

import (
	"fmt"
	"testing"

	"github.com/putrafajarh/bolt/pkg/utils"
	"gotest.tools/v3/assert"
)

func TestCapitalize(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "store",
			expected: "Store",
		},
		{
			input:    "batman",
			expected: "Batman",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "thunder bolt",
			expected: "Thunder bolt",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input: %v, expected: %v", testCase.input, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.Capitalize(testCase.input))
			},
		)
	}
}

func TestSlug(t *testing.T) {
	testCases := []struct {
		input     string
		separator string
		expected  string
	}{
		{
			input:     "lorem ipsum",
			separator: "",
			expected:  "lorem-ipsum",
		},
		{
			input:     "lorem ipsum",
			separator: "_",
			expected:  "lorem_ipsum",
		},
		{
			input:     "  lorem ipsum ",
			separator: "-",
			expected:  "lorem-ipsum",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input: %v, expected: %v", testCase.input, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.Slug(testCase.input, testCase.separator))
			},
		)
	}
}

func TestWords(t *testing.T) {
	testCases := []struct {
		input    string
		limit    uint
		end      string
		expected string
	}{
		{
			input:    "lorem ipsum dolor sit amet",
			limit:    50,
			end:      "",
			expected: "lorem ipsum dolor sit amet",
		},
		{
			input:    "lorem ipsum dolor sit amet",
			limit:    2,
			end:      "...",
			expected: "lorem ipsum...",
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input: %v, limit: %v, end: %v, expected: %v", testCase.input, testCase.limit, testCase.end, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.Words(testCase.input, testCase.limit, testCase.end))
			},
		)
	}
}
