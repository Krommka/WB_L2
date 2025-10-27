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
	"sync"
)

// Processor основная структура утилиты
type Processor struct {
	fetcher       *fetchingResource.Fetcher
	pathConverter *pathConverter.Converter
	cfg           *config.Config
	Assets        map[string]bool
	Links         map[string]bool
	assetsMutex   sync.RWMutex
	linksMutex    sync.RWMutex
}

// New создает процессор
func New(cfg *config.Config) *Processor {
	fetcher := fetchingResource.NewFetcher(10)
	pc := pathConverter.NewPathConverter(cfg.Site)

	resources := make(map[string]bool)
	links := make(map[string]bool)
	return &Processor{fetcher: fetcher, pathConverter: pc, cfg: cfg, Links: links, Assets: resources}
}

// Do стартовый метод процессора
func (processor *Processor) Do() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if isHTML(processor.cfg.Site.String()) {

		processor.Links[processor.cfg.Site.String()] = false
	} else {
		processor.Assets[processor.cfg.Site.String()] = false
	}

	err := processor.process(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (processor *Processor) process(ctx context.Context) error {
	for currentLevel := 0; currentLevel < processor.cfg.Level; currentLevel++ {
		err := processor.processLinks(ctx)
		if err != nil {
			return err
		}
		processor.processAssets(ctx)
	}
	return nil
}

func (processor *Processor) processLinks(ctx context.Context) error {

	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)

	linksToProcess := make([]string, 0)

	for link := range processor.Links {
		if processor.Links[link] == true {
			continue
		}
		linksToProcess = append(linksToProcess, link)
	}

	if len(linksToProcess) == 0 {
		return nil
	}

	for _, link := range linksToProcess {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("process link", "link", link)

			resp, err := processor.fetcher.Fetch(ctx, link)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			parsedUrl, err := url.Parse(link)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			parser := parsingHTML.New(parsedUrl)

			completeHTML, err := parser.Handle(resp.Body)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			processor.assetsMutex.Lock()
			for assetLink := range parser.Assets {
				if _, ok := processor.Assets[assetLink]; !ok {
					processor.Assets[assetLink] = false
				}
			}
			processor.assetsMutex.Unlock()

			processor.linksMutex.Lock()
			for childLink := range parser.Links {
				if _, ok := processor.Links[childLink]; !ok {
					processor.Links[childLink] = false
				}
			}
			processor.linksMutex.Unlock()

			path := processor.pathConverter.URLToLocalPath(parsedUrl)
			err = saveResponse(strings.NewReader(string(completeHTML)), path)
			resp.Body.Close()
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			processor.linksMutex.Lock()
			processor.Links[link] = true
			processor.linksMutex.Unlock()
		}()

	}
	wg.Wait()
	close(errCh)
	return <-errCh
}

func (processor *Processor) processAssets(ctx context.Context) {
	wg := &sync.WaitGroup{}

	assetsToProcess := make([]string, 0)

	for resource := range processor.Assets {
		if processor.Assets[resource] == true {
			continue
		}
		assetsToProcess = append(assetsToProcess, resource)
	}

	if len(assetsToProcess) == 0 {
		return
	}

	for _, resource := range assetsToProcess {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Debug("fetching resource", "resource", resource)
			resp, err := processor.fetcher.Fetch(ctx, resource)
			if err != nil {
				slog.Error("error fetching resource", "resource", resource, "err", err)
				processor.assetsMutex.Lock()
				processor.Assets[resource] = true
				processor.assetsMutex.Unlock()
				return
			}
			parsedUrl, _ := url.Parse(resource)
			path := processor.pathConverter.AssetToLocalPath(parsedUrl)
			err = saveResponse(resp.Body, path)
			if err != nil {
				slog.Error("error saving resource", "resource", resource, "err", err)
			}
			processor.assetsMutex.Lock()
			processor.Assets[resource] = true
			processor.assetsMutex.Unlock()
		}()
	}
	wg.Wait()
	return
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

// saveResponse обрабатывает HTML страницы
func saveResponse(data io.Reader, fileName string) error {
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

	slog.Info("Resource saved", "path", fileName)
	return nil
}
