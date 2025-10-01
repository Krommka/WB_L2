package readinput

import (
	"bufio"
	"os"
)

// ReadInput читает строки из файлов или stdin и добавляет в массив
func ReadInput(files []string) ([]string, error) {
	var lines []string

	if len(files) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		file.Close()

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return lines, nil
}
