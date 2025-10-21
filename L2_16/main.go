package main

import (
	"L2_16/config"
	"L2_16/downloadFile"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	downloader := downloadFile.NewDownloader()

	if err := downloader.Load(cfg); err != nil {
		fmt.Println(err)
	}

}
