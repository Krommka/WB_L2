package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config парсит флаги
type Config struct {
	Level   int
	rawLink string ``
	Site    *url.URL
}

// MustLoad загружает конфиг
func MustLoad() *Config {
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	link, err := createLink(config.rawLink)
	if err != nil {
		panic(err)
	}
	config.Site = link

	return config
}

// ParseFlags парсит входные флаги
func ParseFlags() (*Config, error) {
	config := &Config{}
	if len(os.Args) < 1 {
		return nil, fmt.Errorf("не указаны аргументы")
	}

	fs := flag.NewFlagSet("wget", flag.ContinueOnError)

	fs.IntVar(&config.Level, "l", 1, "level")

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
	config.rawLink = args[0]

	if config.rawLink == "" {
		return nil, fmt.Errorf("сайт не может быть пустым")
	}

	return config, nil
}

// createLink создает начальную ссылку
func createLink(rawLink string) (*url.URL, error) {
	if !strings.Contains(rawLink, "://") {
		rawLink = "https://" + rawLink
	}
	siteLink, err := url.Parse(rawLink)
	if err != nil {
		return nil, err
	}
	if siteLink.Path == "" {
		siteLink.Path = "/"
	}
	return siteLink, nil
}
