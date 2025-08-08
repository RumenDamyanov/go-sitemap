package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.rumenx.com/sitemap"
	chiadapter "go.rumenx.com/sitemap/adapters/chi"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Sitemap routes using adapters
	r.Get("/sitemap.xml", chiadapter.Sitemap(generateSitemap))
	r.Get("/sitemap.txt", chiadapter.SitemapTXT(generateSitemap))
	r.Get("/sitemap.html", chiadapter.SitemapHTML(generateSitemap))

	// Manual handlers for more control
	r.Get("/custom-sitemap.xml", customSitemapHandler)

	// Dynamic sitemap with URL parameters
	r.Get("/sitemap/{category}.xml", dynamicSitemapHandler)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
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

func customSitemapHandler(w http.ResponseWriter, r *http.Request) {
	sm := sitemap.New()

	// Access request for dynamic sitemaps
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		sm.Add("https://example.com/user-agent", time.Now(), 0.3, sitemap.Yearly,
			sitemap.WithTitle("User Agent Specific Page"),
		)
	}

	// Add query-specific content
	category := r.URL.Query().Get("category")
	if category != "" {
		categoryURL := fmt.Sprintf("https://example.com/category/%s", category)
		sm.Add(categoryURL, time.Now(), 0.7, sitemap.Weekly,
			sitemap.WithTitle("Category: "+category),
		)
	}

	xml, err := sm.XML()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xml)
}

func dynamicSitemapHandler(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")

	sm := sitemap.New()

	// Generate category-specific sitemap
	if category != "" {
		// Add category page
		categoryURL := fmt.Sprintf("https://example.com/category/%s", category)
		sm.Add(categoryURL, time.Now(), 0.8, sitemap.Weekly,
			sitemap.WithTitle("Category: "+category),
		)

		// Add sample products in this category
		for i := 1; i <= 5; i++ {
			productURL := fmt.Sprintf("https://example.com/category/%s/product-%d", category, i)
			sm.Add(productURL, time.Now(), 0.7, sitemap.Weekly,
				sitemap.WithTitle(fmt.Sprintf("%s Product %d", category, i)),
			)
		}
	}

	xml, err := sm.XML()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xml)
}
