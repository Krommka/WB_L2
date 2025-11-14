package sorting

import (
	"fmt"
	parsingflags "main/parsingFlags"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func InternalSort(lines []string, config *parsingflags.Config) {
	if len(lines) == 0 {
		return
	}

	if config.CheckSorted {
		if IsSorted(lines, config) {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "sort: disorder in input\n")
			os.Exit(1)
		}
	}

	sorted := SortLines(lines, config)

	for _, line := range sorted {
		fmt.Println(line)
	}
}

// MonthToNumber создает таблицу месяцев
func MonthToNumber(month string) int {
	months := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}
	return months[strings.ToLower(month)]
}

// ParseHumanNumber парсит человеко-читаемые числа (1K, 2M, 3G)
func ParseHumanNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	s = strings.ReplaceAll(s, " ", "")

	multipliers := map[string]float64{
		"K": 1e3, "k": 1e3,
		"M": 1e6, "m": 1e6,
		"G": 1e9, "g": 1e9,
		"T": 1e12, "t": 1e12,
	}

	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(s, suffix) {
			numStr := s[:len(s)-len(suffix)]
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0
			}
			return num * multiplier
		}
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

// GetSortKey создает подстроку с нужным столбцом
func GetSortKey(line string, config *parsingflags.Config) string {
	if config.KeyColumn == 0 {
		if config.IgnoreBlanks {
			return strings.TrimRightFunc(line, unicode.IsSpace)
		}
		return line
	}

	columns := strings.Split(line, config.Delimiter)
	colIndex := config.KeyColumn - 1

	if colIndex < 0 || colIndex >= len(columns) {
		return ""
	}

	value := columns[colIndex]
	if config.IgnoreBlanks {
		value = strings.TrimRightFunc(value, unicode.IsSpace)
	}

	return value
}

// IsSorted проверяет, отсортирован ли массив строк
func IsSorted(lines []string, config *parsingflags.Config) bool {
	for i := 1; i < len(lines); i++ {
		if !CompareStrings(lines[i-1], lines[i], config) {
			return false
		}
	}
	return true
}

// RemoveDuplicates удаляет повторяющиеся строки
func RemoveDuplicates(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

// SortLines сортирует массив строк
func SortLines(lines []string, config *parsingflags.Config) []string {
	sorted := make([]string, len(lines))
	copy(sorted, lines)

	sort.Slice(sorted, func(i, j int) bool {
		return CompareStrings(sorted[i], sorted[j], config)
	})

	if config.Unique {
		sorted = RemoveDuplicates(sorted)
	}

	return sorted
}

// parseNumericValueForSort парсит число из подстроки
func parseNumericValueForSort(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	clean := make([]rune, 0, len(s))
	hasPoint := false

	for i, r := range s {
		if unicode.IsDigit(r) {
			clean = append(clean, r)
		} else if r == '.' && !hasPoint {
			clean = append(clean, r)
			hasPoint = true
		} else if r == '-' && i == 0 {
			clean = append(clean, r)
		} else {
			break
		}
	}

	if len(clean) == 0 {
		return 0, false
	}

	num, err := strconv.ParseFloat(string(clean), 64)
	if err != nil {
		return 0, false
	}

	return num, true
}

// CompareStrings сравнивает две строки согласно конфигу
func CompareStrings(a, b string, config *parsingflags.Config) bool {
	keyA := GetSortKey(a, config)
	keyB := GetSortKey(b, config)
	numberIsSame := false

	if config.Numeric {
		numA, isNumA := parseNumericValueForSort(keyA)
		numB, isNumB := parseNumericValueForSort(keyB)

		if isNumA && isNumB {
			if numA != numB {
				if config.Reverse {
					return numA > numB
				}
				return numA < numB
			}
			numberIsSame = true
		}
		if isNumA && !isNumB {
			return true
		}
		if !isNumA && isNumB {
			return false
		}
	}

	if config.HumanNumeric {
		numA := ParseHumanNumber(keyA)
		numB := ParseHumanNumber(keyB)

		if (numA != 0 || numB != 0) && !(numA == 0 && numB == 0) {
			if config.Reverse {
				return numA > numB
			}
			return numA < numB
		}
	}

	if config.Month {
		monthA := MonthToNumber(keyA)
		monthB := MonthToNumber(keyB)
		if monthA != 0 && monthB != 0 {
			if config.Reverse {
				return monthA > monthB
			}
			return monthA < monthB
		}
	}
	if numberIsSame {
		if config.Reverse {
			return a > b
		}
		return a < b
	}

	if config.Reverse {
		return keyA > keyB
	}
	return keyA < keyB
}
