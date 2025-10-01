package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config хранит конфигурацию grep
type Config struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
}

func ParseFlags() (*Config, []string, error) {
	config := &Config{}

	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	fs.Usage = Usage

	fs.IntVar(&config.After, "A", 0, "print N lines after match")
	fs.IntVar(&config.Before, "B", 0, "print N lines before match")
	fs.IntVar(&config.Context, "C", 0, "print N lines of context")
	fs.BoolVar(&config.Count, "c", false, "print only count")
	fs.BoolVar(&config.IgnoreCase, "i", false, "ignore case")
	fs.BoolVar(&config.Invert, "v", false, "invert match")
	fs.BoolVar(&config.Fixed, "F", false, "fixed string")
	fs.BoolVar(&config.LineNum, "n", false, "print line numbers")

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
					case 'A', 'B', 'C':
						return nil, nil, fmt.Errorf("option -%c cannot be used in combined flags", f)
					case 'c':
						config.Count = true
					case 'i':
						config.IgnoreCase = true
					case 'v':
						config.Invert = true
					case 'F':
						config.Fixed = true
					case 'n':
						config.LineNum = true
					default:
						return nil, nil, fmt.Errorf("unknown option: -%c", f)
					}
				}
			} else {
				switch arg {
				case "-A", "-B", "-C":
					if i+1 < len(args) {
						val := args[i+1]
						n, err := fmt.Sscanf(val, "%d", &config.Context)
						if err != nil || n != 1 {
							return nil, nil, fmt.Errorf("invalid number for %s: %s", arg, val)
						}
						switch arg {
						case "-A":
							config.After = config.Context
						case "-B":
							config.Before = config.Context
						case "-C":
							config.After = config.Context
							config.Before = config.Context
						}
						skipNext = true
					} else {
						return nil, nil, fmt.Errorf("option %s requires an argument", arg)
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

	if len(nonFlagArgs) == 0 {
		return nil, nil, fmt.Errorf("pattern is required")
	}

	config.Pattern = nonFlagArgs[0]
	files := nonFlagArgs[1:]

	if config.After < 0 || config.Before < 0 || config.Context < 0 {
		return nil, nil, fmt.Errorf("context lines count cannot be negative")
	}

	//if config.Context > 0 {
	//	config.After = config.Context
	//	config.Before = config.Context
	//}

	return config, files, nil
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PATTERN [FILE...]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Search for PATTERN in each FILE or standard input.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -A N    print N lines after match\n")
	fmt.Fprintf(os.Stderr, "  -B N    print N lines before match\n")
	fmt.Fprintf(os.Stderr, "  -C N    print N lines of context\n")
	fmt.Fprintf(os.Stderr, "  -c      print only count of matching lines\n")
	fmt.Fprintf(os.Stderr, "  -i      ignore case\n")
	fmt.Fprintf(os.Stderr, "  -v      invert match\n")
	fmt.Fprintf(os.Stderr, "  -F      fixed string (not regexp)\n")
	fmt.Fprintf(os.Stderr, "  -n      print line numbers\n")
}
