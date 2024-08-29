package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		number   int
		expected string
	}{
		{number: 1, expected: "b"},
		{number: 61, expected: "9"},
		{number: 62, expected: "ba"},
		{number: 63, expected: "bb"},
		{number: 125, expected: "cb"},
		{number: 19158, expected: "e9a"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("number=%v", tc.number), func(t *testing.T) {
			actual := Encode(tc.number)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		text     string
		expected int
	}{
		{text: "b", expected: 1},
		{text: "9", expected: 61},
		{text: "ba", expected: 62},
		{text: "bb", expected: 63},
		{text: "cb", expected: 125},
		{text: "e9a", expected: 19158},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("number=%v", tc.text), func(t *testing.T) {
			actual := Decode(tc.text)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestToBase62Digits(t *testing.T) {
	testCases := []struct {
		number   int
		expected []int
	}{
		{number: 1, expected: []int{1}},
		{number: 61, expected: []int{61}},
		{number: 62, expected: []int{1, 0}},
		{number: 63, expected: []int{1, 1}},
		{number: 125, expected: []int{2, 1}},
		{number: 19158, expected: []int{4, 61, 0}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("number=%v", tc.number), func(t *testing.T) {
			actual := toBase62Digits(tc.number)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
