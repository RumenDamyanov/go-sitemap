// Package ginadapter provides Gin framework integration for go-sitemap.
package ginadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.rumenx.com/sitemap"
)

// SitemapGenerator is a function that generates a sitemap.
type SitemapGenerator func() *sitemap.Sitemap

// Sitemap returns a Gin handler that serves a sitemap.
func Sitemap(generator SitemapGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		sm := generator()
		if sm == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		xml, err := sm.XML()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "application/xml")
		c.Data(http.StatusOK, "application/xml", xml)
	}
}

// SitemapTXT returns a Gin handler that serves a sitemap in text format.
func SitemapTXT(generator SitemapGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		sm := generator()
		if sm == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		txt, err := sm.TXT()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/plain")
		c.Data(http.StatusOK, "text/plain", txt)
	}
}

// SitemapHTML returns a Gin handler that serves a sitemap in HTML format.
func SitemapHTML(generator SitemapGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		sm := generator()
		if sm == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		html, err := sm.HTML()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/html")
		c.Data(http.StatusOK, "text/html", html)
	}
}

// SitemapIndex returns a Gin handler that serves a sitemap index.
func SitemapIndex(generator func() *sitemap.Index) gin.HandlerFunc {
	return func(c *gin.Context) {
		idx := generator()
		if idx == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		xml, err := idx.XML()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "application/xml")
		c.Data(http.StatusOK, "application/xml", xml)
	}
}
