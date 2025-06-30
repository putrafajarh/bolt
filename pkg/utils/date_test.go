package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/putrafajarh/bolt/pkg/utils"
	"gotest.tools/v3/assert"
)

func timeLocation(timezone string) *time.Location {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return loc
}

func TestStartOfDay(t *testing.T) {
	testCases := []struct {
		input    time.Time
		timezone string
		expected time.Time
	}{
		{
			input:    time.Unix(-769200739, 0), // Fri Aug 17 1945 05:07:41 GMT+0000
			timezone: "Asia/Jakarta",
			expected: time.Date(1945, time.August, 17, 0, 0, 0, 0, timeLocation("Asia/Jakarta")),
		},
		{
			input:    time.Unix(-769200801, 0), // Fri Aug 17 1945 05:06:39 GMT+0000
			timezone: "UTC",
			expected: time.Date(1945, time.August, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    time.Unix(-769201439, 0), // Fri Aug 17 1945 04:56:01 GMT+0000
			timezone: time.Local.String(),
			expected: time.Date(1945, time.August, 17, 0, 0, 0, 0, time.Local),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input: %v, timezone: %v, expected: %v", testCase.input, testCase.timezone, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected.String(), utils.StartOfDay(testCase.input, testCase.timezone).String())
			},
		)
	}
}

func TestIsSameDay(t *testing.T) {
	testCases := []struct {
		input1   time.Time
		input2   time.Time
		expected bool
	}{
		{ // Sameday different timezone is still the same day
			input1:   time.Date(1945, time.August, 17, 0, 0, 0, 0, timeLocation("Asia/Jakarta")),
			input2:   time.Date(1945, time.August, 17, 0, 0, 0, 0, timeLocation("Asia/Jayapura")),
			expected: true,
		},
		{
			input1:   time.Date(1945, time.August, 17, 0, 0, 0, 0, timeLocation("Asia/Jakarta")),
			input2:   time.Date(1945, time.August, 18, 0, 0, 0, 0, timeLocation("Asia/Jayapura")),
			expected: false,
		},
		{ // Sameday different hour is still the same day
			input1:   time.Date(1945, time.August, 17, 0, 0, 0, 0, timeLocation("Asia/Jakarta")),
			input2:   time.Date(1945, time.August, 17, 12, 0, 0, 0, timeLocation("Brazil/Acre")),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input1: %v, input2: %v, expected: %v", testCase.input1, testCase.input2, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.IsSameDay(testCase.input1, testCase.input2))
			},
		)
	}
}

func TestIsLeapYear(t *testing.T) {
	testCases := []struct {
		input    int
		expected bool
	}{
		{
			input:    2000,
			expected: true,
		},
		{
			input:    2001,
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input: %v, expected: %v", testCase.input, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.IsLeapYear(testCase.input))
			},
		)
	}
}

func TestDiffInDays(t *testing.T) {
	testCases := []struct {
		input1   time.Time
		input2   time.Time
		expected int
	}{
		{
			input1:   time.Now(),
			input2:   time.Now().Add(1 * time.Hour),
			expected: 0,
		},
		{
			input1:   time.Now(),
			input2:   time.Now().Add(-10 * time.Hour),
			expected: 0,
		},
		{
			input1:   time.Now(),
			input2:   time.Now().AddDate(0, 0, 3),
			expected: 3,
		},
		{
			input1:   time.Now(),
			input2:   time.Now().AddDate(0, 0, -2).Add(-1 * time.Millisecond),
			expected: -2,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("TestCase: input1: %v, input2: %v, expected: %v", testCase.input1, testCase.input2, testCase.expected),
			func(t *testing.T) {
				assert.Equal(t, testCase.expected, utils.DiffInDays(testCase.input1, testCase.input2))
			},
		)
	}
}
