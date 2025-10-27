package pathConverter

import (
	"net/url"
	"path/filepath"
	"strings"
)

type Converter struct {
	baseURL *url.URL
}

func NewPathConverter(baseURL *url.URL) *Converter {
	return &Converter{baseURL: baseURL}
}

func (pc *Converter) URLToLocalPath(u *url.URL) string {
	normalized := pc.normalizeURL(u)
	relativePath := pc.getRelativePath(normalized)
	return relativePath

}

func (pc *Converter) ResourceToLocalPath(u *url.URL) string {
	return u.Host + u.Path
}

func (pc *Converter) normalizeURL(u *url.URL) *url.URL {
	normalized := *u
	normalized.Fragment = ""
	normalized.RawQuery = ""
	if normalized.Path == "" {
		normalized.Path = "/"
	}
	return &normalized
}

func (pc *Converter) getRelativePath(u *url.URL) string {

	if u.Host != pc.baseURL.Host {
		return filepath.Join(u.Host, u.Path)
	}
	targetPath := u.Path

	if strings.HasPrefix(targetPath, "/") {
		targetPath = targetPath[1:]
	}

	if targetPath == "" {
		targetPath = "index.html"
	} else if strings.HasSuffix(targetPath, "/") {
		targetPath = targetPath + "index.html"
	}

	return filepath.Join(u.Host, targetPath)

}
