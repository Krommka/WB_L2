package parsingflags

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	level int
	site  string
}

func ParseFlags() (*Config, error) {
	config := Config{}
	if len(os.Args) < 1 {
		return nil, fmt.Errorf("не указаны аргументы")
	}

	fs := flag.NewFlagSet("wget", flag.ContinueOnError)

	fs.IntVar(&config.level, "l", 0, "level")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	if config.level < 0 {
		return nil, fmt.Errorf("уровень должен быть больше или равен 0 (указано: %d)", config.level)
	}

	args := fs.Args()
	if len(args) < 1 {
		return nil, fmt.Errorf("сайт не указан")
	}
	config.site = args[0]

	if config.site == "" {
		return nil, fmt.Errorf("сайт не может быть пустым")
	}

	return &config, nil
}
