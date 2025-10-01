package main

import (
	"fmt"
	"regexp"
	"strings"
)

// MatchResult представляет результат совпадения
type MatchResult struct {
	LineNum int
	Line    string
	Matched bool
}

// Grep выполняет поиск по строкам согласно конфигурации
func Grep(lines []string, config *Config) ([]string, int, error) {
	pattern := config.Pattern
	if config.IgnoreCase {
		if config.Fixed {
			pattern = strings.ToLower(pattern)
		} else {
			pattern = "(?i)" + pattern
		}
	}

	var re *regexp.Regexp
	var err error
	if !config.Fixed {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	matches := make([]MatchResult, len(lines))
	matchCount := 0

	for i, line := range lines {
		var matched bool

		if config.Fixed {
			searchLine := line
			searchPattern := config.Pattern
			if config.IgnoreCase {
				searchLine = strings.ToLower(line)
				searchPattern = strings.ToLower(config.Pattern)
			}
			matched = strings.Contains(searchLine, searchPattern)
		} else {
			if re != nil {
				matched = re.MatchString(line)
			}
		}

		if config.Invert {
			matched = !matched
		}

		matches[i] = MatchResult{
			LineNum: i + 1,
			Line:    line,
			Matched: matched,
		}

		if matched {
			matchCount++
		}
	}

	if config.Count {
		return []string{fmt.Sprintf("%d", matchCount)}, matchCount, nil
	}

	return buildOutput(matches, config), matchCount, nil
}

// buildOutput строит вывод с учетом контекста и номеров строк
func buildOutput(matches []MatchResult, config *Config) []string {
	var result []string
	printed := make(map[int]bool)
	lastMatch := -1

	for i, match := range matches {
		if match.Matched {
			start := i - config.Before
			if start < 0 {
				start = 0
			}
			for j := start; j < i; j++ {
				if !printed[j] {
					if config.Context != 0 && matches[j].LineNum > lastMatch+1 && lastMatch != -1 {
						result = append(result, "--")
					}
					result = appendContextLine(result, matches[j], config, false)
					lastMatch = matches[i].LineNum
					printed[j] = true
				}
			}

			if !printed[i] {
				if config.Context != 0 && matches[i].LineNum > lastMatch+1 && lastMatch != -1 {
					result = append(result, "--")
				}
				result = appendContextLine(result, matches[i], config, true)
				lastMatch = matches[i].LineNum
				printed[i] = true
			}

			end := i + config.After + 1
			if end > len(matches) {
				end = len(matches)
			}
			for j := i + 1; j < end; j++ {
				if !printed[j] {
					result = appendContextLine(result, matches[j], config, false)
					lastMatch = matches[j].LineNum
					printed[j] = true
				}
			}
		}
	}

	return result
}

// appendContextLine добавляет строку в результат с учетом форматирования
func appendContextLine(result []string, match MatchResult, config *Config, isMatch bool) []string {
	line := match.Line

	if config.LineNum {
		prefix := fmt.Sprintf("%d-", match.LineNum)
		if isMatch {
			prefix = fmt.Sprintf("%d:", match.LineNum)
		}
		line = prefix + line
	}

	return append(result, line)
}
