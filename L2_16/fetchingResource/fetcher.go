package fetchingResource

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Fetcher struct {
	semaphore chan struct{}
	client    *http.Client
}

func NewFetcher(concurrency int) *Fetcher {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 3 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}
	return &Fetcher{
		semaphore: make(chan struct{}, concurrency),
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Minute,
		},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (*http.Response, error) {
	f.semaphore <- struct{}{}
	defer func() { <-f.semaphore }()

	ctxTimeout, _ := context.WithTimeout(ctx, time.Minute*10)

	req, err := http.NewRequestWithContext(ctxTimeout, "GET", url, nil)

	if err != nil {
		return nil, err
	}
	slog.Debug("get page", "page", url)

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	//userAgent := fmt.Sprintf("MyGet/1.0 (%s-%s)", runtime.GOOS, runtime.GOARCH)

	req.Header.Set("User-Agent", userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("get url %s, status: %s", url, resp.Status)
	}
	slog.Debug("load success", "link", url)
	return resp, nil
}
