package main

import (
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

var (
	errInvalidString = errors.New("invalid string")
)

func main() {
	//sl := []string{"a4bc2d5e", "abcd", "45", "", "qwe\\4\\5", "qwe\\45", "аб14в2г1д", "\\4\\5", "a10"}
	//for _, v := range sl {
	//
	//	fmt.Printf("Unwrap string: \"%s\". Result: ", v)
	//	str, err := Unpack(v)
	//	if err != nil {
	//		fmt.Printf("%v", err)
	//	}
	//	fmt.Printf("%s\n", str)
	//}
}

func unpack(data string) (string, error) {
	if data == "" {
		return "", nil
	}

	runes := []rune(data)
	var result bytes.Buffer
	var prevRune rune
	escaped := false

	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		if !escaped && currentRune == '\\' {
			escaped = true
			continue
		}

		if escaped {
			result.WriteRune(currentRune)
			prevRune = currentRune
			escaped = false
			continue
		}

		if unicode.IsDigit(currentRune) {
			if prevRune == 0 {
				return "", errInvalidString
			}

			digitStr := string(currentRune)
			for j := i + 1; j < len(runes) && unicode.IsDigit(runes[j]); j++ {
				digitStr += string(runes[j])
				i = j
			}

			count, err := strconv.Atoi(digitStr)
			if err != nil || count <= 0 {
				return "", errInvalidString
			}

			for k := 1; k < count; k++ {
				result.WriteRune(prevRune)
			}

			prevRune = 0
			continue
		}

		result.WriteRune(currentRune)
		prevRune = currentRune
	}

	if escaped {
		return "", errInvalidString
	}

	return result.String(), nil
}
