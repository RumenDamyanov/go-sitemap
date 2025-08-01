// Package echoadapter provides Echo framework integration for go-sitemap.
package echoadapter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-sitemap"
)

// SitemapGenerator is a function that generates a sitemap.
type SitemapGenerator func() *sitemap.Sitemap

// Sitemap returns an Echo handler that serves a sitemap.
func Sitemap(generator SitemapGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		sm := generator()
		if sm == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		xml, err := sm.XML()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Blob(http.StatusOK, "application/xml", xml)
	}
}

// SitemapTXT returns an Echo handler that serves a sitemap in text format.
func SitemapTXT(generator SitemapGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		sm := generator()
		if sm == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		txt, err := sm.TXT()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Blob(http.StatusOK, "text/plain", txt)
	}
}

// SitemapHTML returns an Echo handler that serves a sitemap in HTML format.
func SitemapHTML(generator SitemapGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		sm := generator()
		if sm == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		html, err := sm.HTML()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Blob(http.StatusOK, "text/html", html)
	}
}

// SitemapIndex returns an Echo handler that serves a sitemap index.
func SitemapIndex(generator func() *sitemap.Index) echo.HandlerFunc {
	return func(c echo.Context) error {
		idx := generator()
		if idx == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		xml, err := idx.XML()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Blob(http.StatusOK, "application/xml", xml)
	}
}
