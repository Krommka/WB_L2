package parsingflags

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config описывает структуру флагов
type Config struct {
	KeyColumn    int
	Numeric      bool
	Reverse      bool
	Unique       bool
	Month        bool
	IgnoreBlanks bool
	CheckSorted  bool
	HumanNumeric bool
	Delimiter    string
}

// ParseFlags создает конфиг по объявленным флагам
func ParseFlags() (*Config, []string, error) {
	config := &Config{
		Delimiter: "\t",
	}

	var keyColumn int

	fs := flag.NewFlagSet("sort", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Sort lines of text files\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -k N    sort by column N (1-based)\n")
		fmt.Fprintf(os.Stderr, "  -n      sort numerically\n")
		fmt.Fprintf(os.Stderr, "  -r      reverse sort order\n")
		fmt.Fprintf(os.Stderr, "  -u      output only unique lines\n")
		fmt.Fprintf(os.Stderr, "  -M      sort by month names\n")
		fmt.Fprintf(os.Stderr, "  -b      ignore trailing blanks\n")
		fmt.Fprintf(os.Stderr, "  -c      check if input is sorted\n")
		fmt.Fprintf(os.Stderr, "  -h      sort by human-readable numbers\n")
	}

	fs.IntVar(&keyColumn, "k", 0, "sort by key column number")
	fs.BoolVar(&config.Numeric, "n", false, "sort numerically")
	fs.BoolVar(&config.Reverse, "r", false, "reverse sort order")
	fs.BoolVar(&config.Unique, "u", false, "output only unique lines")
	fs.BoolVar(&config.Month, "M", false, "sort by month names")
	fs.BoolVar(&config.IgnoreBlanks, "b", false, "ignore trailing blanks")
	fs.BoolVar(&config.CheckSorted, "c", false, "check if input is sorted")
	fs.BoolVar(&config.HumanNumeric, "h", false, "sort by human-readable numbers")

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
			if len(arg) > 2 {
				flags := arg[1:]
				for _, f := range flags {
					switch f {
					case 'n':
						config.Numeric = true
					case 'r':
						config.Reverse = true
					case 'u':
						config.Unique = true
					case 'M':
						config.Month = true
					case 'b':
						config.IgnoreBlanks = true
					case 'c':
						config.CheckSorted = true
					case 'h':
						config.HumanNumeric = true
					case 'k':
						if i+1 < len(args) {
							val := args[i+1]
							if n, err := strconv.Atoi(val); err == nil && n >= 0 {
								keyColumn = n
								skipNext = true
							} else {
								return nil, nil, fmt.Errorf("invalid column number: %s", val)
							}
						} else {
							return nil, nil, fmt.Errorf("option -k requires an argument")
						}
					default:
						return nil, nil, fmt.Errorf("unknown option: -%c", f)
					}
				}
			} else {
				switch arg {
				case "-k":
					if i+1 < len(args) {
						val := args[i+1]
						if n, err := strconv.Atoi(val); err == nil && n >= 0 {
							keyColumn = n
							skipNext = true
						} else {
							return nil, nil, fmt.Errorf("invalid column number: %s", val)
						}
					} else {
						return nil, nil, fmt.Errorf("option -k requires an argument")
					}
				default:
					if err := fs.Parse([]string{arg}); err != nil {
						return nil, nil, err
					}
				}
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}

	config.KeyColumn = keyColumn

	if config.Numeric && config.Month {
		return nil, nil, fmt.Errorf("conflicting options: -n and -M")
	}
	if config.Numeric && config.HumanNumeric {
		return nil, nil, fmt.Errorf("conflicting options: -n and -h")
	}
	if config.Month && config.HumanNumeric {
		return nil, nil, fmt.Errorf("conflicting options: -M and -h")
	}

	return config, nonFlagArgs, nil
}

// Usage создает комментарии для терминала
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE...]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Sort lines of text files\n\n")
	fmt.Fprintf(os.Stderr, "Mandatory options (like GNU sort):\n")
	fmt.Fprintf(os.Stderr, "  -k N          sort by column N (1-based)\n")
	fmt.Fprintf(os.Stderr, "  -n            sort numerically\n")
	fmt.Fprintf(os.Stderr, "  -r            reverse sort order\n")
	fmt.Fprintf(os.Stderr, "  -u            output only unique lines\n")
	fmt.Fprintf(os.Stderr, "\nAdditional options:\n")
	fmt.Fprintf(os.Stderr, "  -M            sort by month names\n")
	fmt.Fprintf(os.Stderr, "  -b            ignore trailing blanks\n")
	fmt.Fprintf(os.Stderr, "  -c            check if input is sorted\n")
	fmt.Fprintf(os.Stderr, "  -h            sort by human-readable numbers\n")
	fmt.Fprintf(os.Stderr, "\nCombined flags are supported: -nr, -nru, etc.\n")
}
