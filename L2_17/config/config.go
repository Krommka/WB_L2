package config

import (
	"errors"
	"flag"
	"time"
)

// Config представляет конфигурацию telnet-клиента
type Config struct {
	Host    string
	Port    string
	Timeout time.Duration
}

// MustLoad загружает конфигурацию или паникует при ошибке
func MustLoad() *Config {
	cfg, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	return cfg
}

// ParseFlags парсит аргументы командной строки
func ParseFlags() (*Config, error) {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", time.Second*10,
		"connetion timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return nil, errors.New("usage: mytelnet [--timeout=10s] host port")
	}

	host := args[0]
	port := args[1]

	return &Config{
		host,
		port,
		timeout,
	}, nil
}
