package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Level     int
	Site      string
	Resources []string
	Links     []Page
}

type Page struct {
	address string
	level   int
}

func MustLoad() *Config {
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	config.Resources = make([]string, 0)
	config.Links = make([]Page, 0)

	return config
}

func ParseFlags() (*Config, error) {
	config := Config{}
	if len(os.Args) < 1 {
		return nil, fmt.Errorf("не указаны аргументы")
	}

	fs := flag.NewFlagSet("wget", flag.ContinueOnError)

	fs.IntVar(&config.Level, "l", 0, "level")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	if config.Level < 0 {
		return nil, fmt.Errorf("уровень должен быть больше или равен 0 (указано: %d)", config.Level)
	}

	args := fs.Args()
	if len(args) < 1 {
		return nil, fmt.Errorf("сайт не указан")
	}
	config.Site = args[0]

	if config.Site == "" {
		return nil, fmt.Errorf("сайт не может быть пустым")
	}

	return &config, nil
}
