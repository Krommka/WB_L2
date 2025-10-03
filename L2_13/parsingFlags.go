package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	Start, End int
}

type Config struct {
	Delimiter string
	Separated bool
	Ranges    []Range
}

func ParseFlags() (*Config, []string, error) {
	config := &Config{
		Delimiter: "\t",
	}

	fs := flag.NewFlagSet("cut", flag.ContinueOnError)
	fs.Usage = func() {
		Usage()
	}

	var fieldsString string

	fs.StringVar(&fieldsString, "f", "", "print chosen FIELDS in data")
	fs.StringVar(&config.Delimiter, "d", "\t", "use chosen DELIMITER as separator")
	fs.BoolVar(&config.Separated, "s", false, "print ONLY line with delimiter")

	// Обрабатываем аргументы вручную для поддержки комбинированных флагов
	args := os.Args[1:]
	var nonFlagArgs []string
	var skipNext bool

	for i := 0; i < len(args); i++ {
		if skipNext {
			skipNext = false
			continue
		}

		arg := args[i]

		if strings.HasPrefix(arg, "-") && len(arg) > 1 && !strings.HasPrefix(arg, "--") {
			// Комбинированные флаги
			if len(arg) > 2 {
				flags := arg[1:]
				for _, f := range flags {
					switch f {
					case 'f', 'd':
						// Флаги с аргументами в комбинации не поддерживаем
						return nil, nil, fmt.Errorf("option -%c cannot be used in combined flags", f)
					case 's':
						config.Separated = true
					default:
						return nil, nil, fmt.Errorf("unknown option: -%c", f)
					}
				}
			} else {
				// Одиночные флаги
				switch arg {
				case "-f", "-d":
					if i+1 < len(args) {
						// Временно устанавливаем для парсинга
						tempArgs := []string{arg, args[i+1]}
						if err := fs.Parse(tempArgs); err != nil {
							return nil, nil, err
						}
						skipNext = true
					} else {
						return nil, nil, fmt.Errorf("option %s requires an argument", arg)
					}
				default:
					// Обрабатываем через flagset
					if err := fs.Parse([]string{arg}); err != nil {
						return nil, nil, err
					}
				}
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}

	// Если fieldsString не установлен через флаг, берем первый не-флаговый аргумент
	if fieldsString == "" && len(nonFlagArgs) > 0 {
		fieldsString = nonFlagArgs[0]
		nonFlagArgs = nonFlagArgs[1:]
	}

	if fieldsString == "" {
		return nil, nil, fmt.Errorf("fields not specified")
	}

	ranges, err := parseRangeString(fieldsString)
	if err != nil {
		return nil, nil, err
	}
	config.Ranges = ranges

	return config, nonFlagArgs, nil
}

func parseRangeString(s string) ([]Range, error) {
	ranges := make([]Range, 0)

	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid field range: %s", part)
			}

			startStr := strings.TrimSpace(rangeParts[0])
			endStr := strings.TrimSpace(rangeParts[1])

			if startStr == "" || endStr == "" {
				return nil, fmt.Errorf("invalid field range: %s", part)
			}

			start, err := strconv.Atoi(startStr)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", startStr)
			}

			end, err := strconv.Atoi(endStr)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", endStr)
			}

			if start < 1 || end < 1 {
				return nil, fmt.Errorf("fields are numbered from 1")
			}

			if start > end {
				return nil, fmt.Errorf("invalid decreasing range")
			}

			ranges = append(ranges, Range{Start: start, End: end})
		} else {
			value, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			if value < 1 {
				return nil, fmt.Errorf("fields are numbered from 1")
			}
			ranges = append(ranges, Range{Start: value, End: value})
		}
	}

	return ranges, nil
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -f LIST [OPTIONS] [FILE...]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Print selected parts of lines from each FILE to standard output.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -f LIST   select only these fields; also print any line that contains no delimiter character unless the -s option is specified\n")
	fmt.Fprintf(os.Stderr, "  -d DELIM  use DELIM instead of TAB for field delimiter\n")
	fmt.Fprintf(os.Stderr, "  -s        do not print lines not containing delimiters\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s -f 1,3-5 file.txt    # output 1st and 3rd to 5th fields\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -f 2- -d ',' file.csv # output from 2nd field to end with comma delimiter\n", os.Args[0])
}
