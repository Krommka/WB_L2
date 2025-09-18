package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpack(t *testing.T) {

	testCases := []struct {
		name     string
		data     string
		expected string
		hasError bool
	}{
		{"Basic unpack", "a4bc2d5e", "aaaabccddddde", false},
		{"No numbers", "abcd", "abcd", false},
		{"Only numbers", "45", "", true},
		{"Empty string", "", "", false},
		{"Escaped numbers", "qwe\\4\\5", "qwe45", false},
		{"Escaped with number", "qwe\\45", "qwe44444", false},
		{"Cyrillic", "аб4в2г1д", "аббббввгд", false},
		{"Single char", "*", "*", false},
		{"Starts with number", "2asdfr", "", true},
		{"Zero count", "asd0fr", "", true},
		{"Large number", "a10", "aaaaaaaaaa", false},
		{"Multiple digits", "a12b3", "aaaaaaaaaaaabbb", false},
		{"Invalid escape end", "qwe\\", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := unpack(tc.data)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
