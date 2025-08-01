// Package chiadapter provides Chi router integration for go-sitemap.
package chiadapter

import (
	"net/http"

	"github.com/rumendamyanov/go-sitemap"
)

// SitemapGenerator is a function that generates a sitemap.
type SitemapGenerator func() *sitemap.Sitemap

// Sitemap returns an HTTP handler that serves a sitemap.
func Sitemap(generator SitemapGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := generator()
		if sm == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		xml, err := sm.XML()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml)
	}
}

// SitemapTXT returns an HTTP handler that serves a sitemap in text format.
func SitemapTXT(generator SitemapGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := generator()
		if sm == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		txt, err := sm.TXT()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(txt)
	}
}

// SitemapHTML returns an HTTP handler that serves a sitemap in HTML format.
func SitemapHTML(generator SitemapGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := generator()
		if sm == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		html, err := sm.HTML()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	}
}

// SitemapIndex returns an HTTP handler that serves a sitemap index.
func SitemapIndex(generator func() *sitemap.Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idx := generator()
		if idx == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		xml, err := idx.XML()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml)
	}
}
