package downloadFile

import (
	config "L2_16/config"
	"net/http"
	"time"
)

type DownLoader struct {
	Config config.Config
	client *http.Client
}

func NewDownloader() DownLoader {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return DownLoader{client: client}
}
