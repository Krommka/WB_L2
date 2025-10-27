package processing

import (
	"L2_16/config"
	"L2_16/fetchingResource"
	"L2_16/parsingHTML"
	"L2_16/pathConverter"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Processor struct {
	fetcher       *fetchingResource.Fetcher
	pathConverter *pathConverter.Converter
	cfg           *config.Config
	Assets        map[*url.URL]bool
	Links         map[*url.URL]bool
	Downloaded    map[string]bool
}

func New(cfg *config.Config) *Processor {
	fetcher := fetchingResource.NewFetcher(10)
	pc := pathConverter.NewPathConverter(cfg.Site)

	resources := make(map[*url.URL]bool)
	links := make(map[*url.URL]bool)
	downloaded := make(map[string]bool)
	return &Processor{fetcher: fetcher, pathConverter: pc, cfg: cfg, Links: links, Assets: resources,
		Downloaded: downloaded}
}

func (processor *Processor) Do() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if isHTML(processor.cfg.Site.String()) {

		processor.Links[processor.cfg.Site] = false
	} else {
		processor.Assets[processor.cfg.Site] = false
	}

	err := processor.process(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (processor *Processor) process(ctx context.Context) error {
	for currentLevel := 0; currentLevel < processor.cfg.Level; currentLevel++ {
		for link := range processor.Links {
			if processor.Links[link] == true {
				continue
			}
			slog.Info("process link", "link", link)
			resp, err := processor.fetcher.Fetch(ctx, link.String())
			if err != nil {
				return err
			}
			parser := parsingHTML.New(link)

			completeHTML, err := parser.Handle(resp.Body)
			if err != nil {
				return err
			}

			for assetLink := range parser.Assets {
				if _, ok := processor.Assets[assetLink]; !ok {
					processor.Assets[assetLink] = false
				}
			}
			for childLink := range parser.Links {
				if _, ok := processor.Links[childLink]; !ok {
					processor.Links[childLink] = false
				}
			}
			path := processor.pathConverter.URLToLocalPath(link)
			err = SaveResponse(strings.NewReader(string(completeHTML)), path)
			resp.Body.Close()
			if err != nil {
				return err
			}
			processor.Links[link] = true
		}

		for resource := range processor.Assets {
			if processor.Assets[resource] == true {
				continue
			}
			slog.Debug("fetching resource", "resource", resource)
			resp, err := processor.fetcher.Fetch(ctx, resource.String())
			if err != nil {
				slog.Error("error fetching resource", "resource", resource, "err", err)
				processor.Assets[resource] = true
				continue
			}
			path := processor.pathConverter.ResourceToLocalPath(resource)
			err = SaveResponse(resp.Body, path)
			if err != nil {
				slog.Error("error saving resource", "resource", resource, "err", err)
			}
			processor.Assets[resource] = true
		}
	}
	return nil
}

func isHTML(link string) bool {
	path := strings.ToLower(link)
	ext := filepath.Ext(path)

	htmlExtensions := map[string]bool{
		".html": true, ".htm": true, ".xhtml": true,
		".php": true, ".asp": true, ".aspx": true, ".jsp": true,
		"": true,
	}
	return htmlExtensions[ext]
}

// SaveResponse обрабатывает HTML страницы
func SaveResponse(data io.Reader, fileName string) error {
	slog.Debug("saving resource", "path", fileName)
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	slog.Debug("Resource saved", "path", fileName)
	return nil
}
