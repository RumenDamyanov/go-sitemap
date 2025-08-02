package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rumendamyanov/go-sitemap"
	fiberadapter "github.com/rumendamyanov/go-sitemap/adapters/fiber"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Sitemap Example v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Sitemap routes using adapters
	app.Get("/sitemap.xml", fiberadapter.Sitemap(generateSitemap))
	app.Get("/sitemap.txt", fiberadapter.SitemapTXT(generateSitemap))
	app.Get("/sitemap.html", fiberadapter.SitemapHTML(generateSitemap))

	// Manual handlers for more control
	app.Get("/custom-sitemap.xml", customSitemapHandler)

	fmt.Println("Starting server on :8080")
	app.Listen(":8080")
}

func generateSitemap() *sitemap.Sitemap {
	sm := sitemap.New()

	// Add homepage
	sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily,
		sitemap.WithTitle("Homepage"),
	)

	// Add API endpoints
	apiEndpoints := []string{
		"/api/users",
		"/api/products",
		"/api/orders",
	}

	for _, endpoint := range apiEndpoints {
		sm.Add("https://example.com"+endpoint, time.Now(), 0.5, sitemap.Weekly,
			sitemap.WithTitle("API Endpoint: "+endpoint),
		)
	}

	// Add product pages with images
	products := []struct {
		id    int
		name  string
		image string
	}{
		{1, "Awesome Product", "/images/product-1.jpg"},
		{2, "Great Product", "/images/product-2.jpg"},
		{3, "Best Product", "/images/product-3.jpg"},
	}

	for _, product := range products {
		images := []sitemap.Image{
			{
				URL:   fmt.Sprintf("https://example.com%s", product.image),
				Title: product.name,
			},
		}

		productURL := fmt.Sprintf("https://example.com/products/%d", product.id)
		sm.Add(productURL, time.Now(), 0.8, sitemap.Weekly,
			sitemap.WithTitle(product.name),
			sitemap.WithImages(images),
		)
	}

	// Add blog posts
	for i := 1; i <= 10; i++ {
		blogURL := fmt.Sprintf("https://example.com/blog/post-%d", i)
		sm.Add(blogURL, time.Now().AddDate(0, 0, -i), 0.6, sitemap.Monthly,
			sitemap.WithTitle(fmt.Sprintf("Blog Post %d", i)),
		)
	}

	return sm
}

func customSitemapHandler(c *fiber.Ctx) error {
	sm := sitemap.New()

	// Access Fiber context for dynamic sitemaps
	userAgent := c.Get("User-Agent")
	if userAgent != "" {
		sm.Add("https://example.com/user-agent", time.Now(), 0.3, sitemap.Yearly,
			sitemap.WithTitle("User Agent Specific Page"),
		)
	}

	// Add query-specific content
	category := c.Query("category")
	if category != "" {
		categoryURL := fmt.Sprintf("https://example.com/category/%s", category)
		sm.Add(categoryURL, time.Now(), 0.7, sitemap.Weekly,
			sitemap.WithTitle("Category: "+category),
		)
	}

	xml, err := sm.XML()
	if err != nil {
		return c.SendStatus(500)
	}

	c.Set("Content-Type", "application/xml")
	return c.Send(xml)
}
