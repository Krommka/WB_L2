package main

import (
	"strings"
)

// Cut обрабатывает строки согласно конфигурации
func Cut(lines []string, config *Config) []string {
	res := make([]string, 0)
	for _, line := range lines {
		fields := strings.Split(line, config.Delimiter)

		if config.Separated && len(fields) == 1 {
			continue
		}

		selectedFields := selectFields(fields, config.Ranges)

		if len(selectedFields) > 0 {
			res = append(res, strings.Join(selectedFields, config.Delimiter))
			continue
		}
		if len(fields) == 1 {
			res = append(res, line)
			continue
		}
		res = append(res, "")
	}
	return res
}

// selectFields выбирает поля согласно диапазонам
func selectFields(fields []string, ranges []Range) []string {
	var result []string

	if len(ranges) > 0 {
		// Собираем все номера полей в правильном порядке
		fieldOrder := getFieldOrder(ranges)

		for _, fieldNum := range fieldOrder {
			idx := fieldNum - 1
			if idx >= 0 && idx < len(fields) {
				result = append(result, fields[idx])
			}
		}
	}

	return result
}

// getFieldOrder возвращает упорядоченный список полей согласно диапазонам
func getFieldOrder(ranges []Range) []int {
	var order []int
	added := make(map[int]bool)

	for _, r := range ranges {
		for i := r.Start; i <= r.End; i++ {
			if !added[i] {
				order = append(order, i)
				added[i] = true
			}
		}
	}

	return order
}
