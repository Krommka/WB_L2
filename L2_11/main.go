package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	s1 := []string{"пятак", "пятка", "тяпка", "листок", "столик", "слиток", "стол"}
	s2 := []string{"пятак", "пятка", "тяпка", "столик", "листок", "слиток", "стол"}
	s3 := []string{"пятак", "пятка", "Тяпка", "столик", "листок", "слиток", "стол"}
	s4 := []string{}
	s5 := []string{"1234", "2341", "3334", "4444", "54444", "6565", "5656"}
	sl := [][]string{s1, s2, s3, s4, s5}
	for i, s := range sl {
		fmt.Printf("s%d\n", i)
		mp := FindAnagrams(s)
		for k, _ := range mp {
			fmt.Printf("- \"%s\": [\"%s\"]\n", k, strings.Join(mp[k], "\", \""))
		}
		fmt.Printf("\n")
	}
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

		groups[key] = append(groups[key], lower)
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
