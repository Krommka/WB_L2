package parsingHTML

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
)

// HTMLHandler обработчик HTML
type HTMLHandler struct {
	link   *url.URL
	base   *url.URL
	Links  map[*url.URL]bool
	Assets map[*url.URL]bool
}

// New создает парсер HTML
func New(resource *url.URL) *HTMLHandler {
	links := make(map[*url.URL]bool)
	assets := make(map[*url.URL]bool)
	return &HTMLHandler{
		link:   resource,
		base:   nil,
		Links:  links,
		Assets: assets,
	}
}

// Handle обрабатывает HTML собирая все ссылки и трансформируя пути
func (handler *HTMLHandler) Handle(data io.ReadCloser) ([]byte, error) {
	bytes, err := io.ReadAll(data)
	defer data.Close()
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	doc, err := html.Parse(strings.NewReader(string(bytes)))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}
	handler.setupBaseHref(doc)
	if handler.base == nil {
		handler.base = handler.link
	}
	handler.extractLinksAndResources(doc)

	transformedHTML, err := handler.transformHTMLLinks(doc)
	if err != nil {
		return nil, fmt.Errorf("transform HTML Links: %w", err)
	}

	return transformedHTML, nil

}

// setupBaseHref Устанавливает базовую ссылку, от которой считаются относительные пути
func (handler *HTMLHandler) setupBaseHref(n *html.Node) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "base":
			if href := getAttribute(n, "href"); href != "" {
				if absoluteURL := toAbsoluteURL(href, handler.base); absoluteURL != nil {
					handler.base = absoluteURL
					slog.Debug("Found base", "base", absoluteURL.String())
				}
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		handler.setupBaseHref(child)
	}
}

// extractLinksAndResources извлекает ресурсы любых доменов, но ссылки только своего домена
func (handler *HTMLHandler) extractLinksAndResources(n *html.Node) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			if href := getAttribute(n, "href"); href != "" {
				if absoluteURL := toAbsoluteURL(href, handler.base); absoluteURL != nil {
					if isSameDomain(absoluteURL, handler.link) {
						if _, ok := handler.Links[absoluteURL]; !ok {
							handler.Links[absoluteURL] = true
							slog.Info("Found resource", "type", n.Data, "link", absoluteURL.String())
						}
					}
				}
			}

		case "img":
			if src := getAttribute(n, "src"); src != "" {
				if absoluteURL := toAbsoluteURL(src, handler.base); absoluteURL != nil {
					if _, ok := handler.Assets[absoluteURL]; !ok {
						handler.Assets[absoluteURL] = true
						slog.Info("Found resource", "type", n.Data, "link", absoluteURL.String())
					}
				}
			}

		case "link":
			if rel := getAttribute(n, "rel"); rel == "stylesheet" || rel == "icon" || rel == "shortcut icon" {
				if href := getAttribute(n, "href"); href != "" {
					if absoluteURL := toAbsoluteURL(href, handler.base); absoluteURL != nil {
						if _, ok := handler.Assets[absoluteURL]; !ok {
							handler.Assets[absoluteURL] = true
							slog.Info("Found resource", "type", n.Data, "link", absoluteURL.String())
						}
					}
				}
			}

		case "script":
			if src := getAttribute(n, "src"); src != "" {
				if absoluteURL := toAbsoluteURL(src, handler.base); absoluteURL != nil {
					if _, ok := handler.Assets[absoluteURL]; !ok {
						handler.Assets[absoluteURL] = true
						slog.Info("Found resource", "type", n.Data, "link", absoluteURL.String())
					}
				}
			}

		case "iframe", "embed", "source", "track", "audio", "video":
			if src := getAttribute(n, "src"); src != "" {
				if absoluteURL := toAbsoluteURL(src, handler.base); absoluteURL != nil {
					if _, ok := handler.Assets[absoluteURL]; !ok {
						handler.Assets[absoluteURL] = true
						slog.Info("Found resource", "type", n.Data, "link", absoluteURL.String())
					}
				}
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		handler.extractLinksAndResources(child)
	}
}

// transformHTMLLinks преобразует только ресурсы для отображения в относительные пути
func (handler *HTMLHandler) transformHTMLLinks(doc *html.Node) ([]byte, error) {
	var transformAttributes func(*html.Node)
	transformAttributes = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "img", "link", "script", "iframe", "embed", "source", "track", "audio", "video":
				if src := getAttribute(n, "src"); src != "" {
					if absoluteURL := toAbsoluteURL(src, handler.base); absoluteURL != nil {
						relativePath := toRelativePath(absoluteURL, handler.link)
						setAttribute(n, "src", relativePath)
					}
				}
				if href := getAttribute(n, "href"); href != "" {
					if absoluteURL := toAbsoluteURL(href, handler.base); absoluteURL != nil {
						relativePath := toRelativePath(absoluteURL, handler.link)
						setAttribute(n, "href", relativePath)
					}
				}
			case "a":
				if href := getAttribute(n, "href"); href != "" {
					slog.Debug("Found href: ", href)
					if absoluteURL := toAbsoluteURL(href, handler.base); absoluteURL != nil {
						if isSameDomain(absoluteURL, handler.link) {
							relativePath := toRelativePath(absoluteURL, handler.link)
							setAttribute(n, "href", relativePath+"/index.html")
						}
					}
				}
			case "base":
				if href := getAttribute(n, "href"); href != "" {
					setAttribute(n, "href", "")
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			transformAttributes(child)
		}
	}

	transformAttributes(doc)

	var buf strings.Builder
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

// getAttribute возвращает значение атрибута HTML элемента
func getAttribute(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

// setAttribute устанавливает значение атрибута HTML элемента
func setAttribute(n *html.Node, attrName, value string) {
	for i, attr := range n.Attr {
		if attr.Key == attrName {
			n.Attr[i].Val = value
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: attrName, Val: value})
}

// toRelativePath преобразует абсолютный URL в относительный путь для локального хранения
func toRelativePath(link, baseURL *url.URL) string {
	fromLocal := filepath.Join(baseURL.Hostname(), baseURL.Path)
	toLocal := filepath.Join(link.Hostname(), link.Path)

	fromDir := filepath.Dir(fromLocal)

	if isSameDomain(link, baseURL) {
		relative, err := filepath.Rel(filepath.Dir(baseURL.Path), link.Path)
		if err != nil {
			return link.Path
		}
		return filepath.ToSlash(relative)
	}

	relative, err := filepath.Rel(fromDir, toLocal)
	if err != nil {
		return toLocal
	}

	return "../" + filepath.ToSlash(relative)
}

// isSameDomain проверяет, принадлежит ли URL тому же домену
func isSameDomain(urlStr, baseDomain *url.URL) bool {
	return urlStr.Hostname() == baseDomain.Hostname()
}

// toAbsoluteURL преобразует относительный URL в абсолютный
func toAbsoluteURL(link string, baseURL *url.URL) *url.URL {
	if link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
		return nil
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return nil
	}

	if parsed.IsAbs() {
		return parsed
	}

	return baseURL.ResolveReference(parsed)
}
