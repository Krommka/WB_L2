package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	sl := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	//sl := []string{"fdg", "asdsf", "ds", "dfg", "wsad", "sdfgg", "sadf"}
	mp := FindAnagrams(sl)
	fmt.Println(mp)
}

// FindAnagrams находит все множества анаграмм по заданному словарю
func FindAnagrams(words []string) map[string][]string {
	groups := make(map[string][]string)
	seen := make(map[string]bool)

	for _, word := range words {
		lower := strings.ToLower(word)
		if seen[lower] {
			continue
		}

		runes := []rune(lower)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		key := string(runes)

		groups[key] = append(groups[key], word)
	}

	result := make(map[string][]string)
	for _, group := range groups {
		if len(group) > 1 {
			key := group[0]
			sort.Strings(group)
			result[key] = group
		}
	}

	return result
}
