package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rumendamyanov/go-sitemap"
	ginadapter "github.com/rumendamyanov/go-sitemap/adapters/gin"
)

func main() {
	r := gin.Default()

	// Sitemap routes using adapters
	r.GET("/sitemap.xml", ginadapter.Sitemap(generateSitemap))
	r.GET("/sitemap.txt", ginadapter.SitemapTXT(generateSitemap))
	r.GET("/sitemap.html", ginadapter.SitemapHTML(generateSitemap))

	// Manual handlers for more control
	r.GET("/custom-sitemap.xml", customSitemapHandler)

	r.Run(":8080")
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
				URL:     "https://example.com" + product.image,
				Title:   product.name,
				Caption: "Product image for " + product.name,
			},
		}

		sm.Add("https://example.com/products/"+fmt.Sprintf("%d", product.id), time.Now(), 0.8, sitemap.Weekly,
			sitemap.WithTitle(product.name),
			sitemap.WithImages(images),
		)
	}

	return sm
}

func customSitemapHandler(c *gin.Context) {
	sm := sitemap.New()

	// Add custom URLs based on request context
	userAgent := c.GetHeader("User-Agent")
	if userAgent != "" {
		sm.Add("https://example.com/user-agent-specific", time.Now(), 0.5, sitemap.Daily,
			sitemap.WithTitle("User Agent Specific Page"),
		)
	}

	// Add query parameter specific URLs
	if c.Query("mobile") == "true" {
		sm.Add("https://example.com/mobile", time.Now(), 0.7, sitemap.Weekly,
			sitemap.WithTitle("Mobile Version"),
		)
	}

	xml, err := sm.XML()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.Header("Content-Type", "application/xml")
	c.Data(200, "application/xml", xml)
}
