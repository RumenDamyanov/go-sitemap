// Package fiberadapter provides Fiber framework integration for go-sitemap.
package fiberadapter

import (
	"github.com/gofiber/fiber/v2"
	"go.rumenx.com/sitemap"
)

// SitemapGenerator is a function that generates a sitemap.
type SitemapGenerator func() *sitemap.Sitemap

// Sitemap returns a Fiber handler that serves a sitemap.
func Sitemap(generator SitemapGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sm := generator()
		if sm == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		xml, err := sm.XML()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "application/xml")
		return c.Send(xml)
	}
}

// SitemapTXT returns a Fiber handler that serves a sitemap in text format.
func SitemapTXT(generator SitemapGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sm := generator()
		if sm == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		txt, err := sm.TXT()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "text/plain")
		return c.Send(txt)
	}
}

// SitemapHTML returns a Fiber handler that serves a sitemap in HTML format.
func SitemapHTML(generator SitemapGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sm := generator()
		if sm == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		html, err := sm.HTML()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "text/html")
		return c.Send(html)
	}
}

// SitemapIndex returns a Fiber handler that serves a sitemap index.
func SitemapIndex(generator func() *sitemap.Index) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idx := generator()
		if idx == nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		xml, err := idx.XML()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "application/xml")
		return c.Send(xml)
	}
}
