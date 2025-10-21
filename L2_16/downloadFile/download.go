package downloadFile

import (
	"L2_16/config"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (d *DownLoader) Load(cfg *config.Config) error {
	rawLink := cfg.Site

	link, err := createLink(rawLink)
	if err != nil {
		return err
	}

	resp, err := d.getPage(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = saveResponse(resp, link); err != nil {
		return err
	}
	return nil

}

func createLink(rawLink string) (*url.URL, error) {
	if !strings.Contains(rawLink, "://") {
		rawLink = "https://" + rawLink
	}
	siteLink, err := url.Parse(rawLink)
	if err != nil {
		return nil, err
	}
	return siteLink, nil
}

func (d *DownLoader) getPage(link *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("request ", req)

	userAgent := fmt.Sprintf("MyWGet/1.0 (%s-%s)", runtime.GOOS, runtime.GOARCH)

	req.Header.Set("User-Agent", userAgent)
	fmt.Println("agent ", userAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp ", resp)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	fmt.Println("load success")
	return resp, nil
}

func saveResponse(resp *http.Response, link *url.URL) error {
	contentType := resp.Header.Get("Content-Type")
	isHTML := strings.Contains(contentType, "text/html")

	filename := generateFilename(link, isHTML)

	if isHTML {
		if err := savePage(resp.Body, filename); err != nil {
			return err
		}
	} else {
		if err := saveElement(resp.Body, filename); err != nil {
			return err
		}
	}
	return nil
}

func generateFilename(link *url.URL, isHTML bool) string {
	if isHTML {
		path := strings.Trim(link.Path, "/")
		if path == "" {
			return "index.html"
		}
		if !strings.Contains(path, ".") {
			return path + ".html"
		}
		return path
	}
	path := link.Host + link.Path
	return strings.TrimPrefix(path, "/")
}

func saveElement(data io.Reader, fileName string) error {
	fmt.Println("start load element ", fileName)
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	fmt.Printf("Element saved: %s\n", fileName)
	return nil
}

func savePage(data io.Reader, fileName string) error {
	fmt.Println("start load page ", fileName)

	dir := filepath.Dir(fileName)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directories: %w", err)
		}
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

	fmt.Printf("Page saved: %s\n", fileName)
	return nil
}
