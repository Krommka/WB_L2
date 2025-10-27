package main

import (
	"L2_16/config"
	"L2_16/processing"

	"fmt"
	"log/slog"
	"os"
)

func init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}

func main() {
	cfg := config.MustLoad()

	processor := processing.New(cfg)

	if err := processor.Do(); err != nil {
		fmt.Println(err)
	}

}
