package main

import (
	"log"
	"net/http"
	"time"

	"go.rumenx.com/sitemap"
)

func main() {
	http.HandleFunc("/sitemap.xml", sitemapHandler)
	http.HandleFunc("/sitemap.txt", sitemapTxtHandler)
	http.HandleFunc("/sitemap.html", sitemapHtmlHandler)

	log.Println("Server starting on :8080")
	log.Println("Visit http://localhost:8080/sitemap.xml")
	log.Println("Visit http://localhost:8080/sitemap.txt")
	log.Println("Visit http://localhost:8080/sitemap.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	sm := generateSitemap()

	w.Header().Set("Content-Type", "application/xml")
	xml, err := sm.XML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(xml)
}

func sitemapTxtHandler(w http.ResponseWriter, r *http.Request) {
	sm := generateSitemap()

	w.Header().Set("Content-Type", "text/plain")
	txt, err := sm.TXT()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(txt)
}

func sitemapHtmlHandler(w http.ResponseWriter, r *http.Request) {
	sm := generateSitemap()

	w.Header().Set("Content-Type", "text/html")
	html, err := sm.HTML()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(html)
}

func generateSitemap() *sitemap.Sitemap {
	sm := sitemap.New()

	// Add homepage
	sm.Add("http://localhost:8080/", time.Now(), 1.0, sitemap.Daily,
		sitemap.WithTitle("Homepage"),
	)

	// Add some pages with different priorities and change frequencies
	pages := []struct {
		url        string
		title      string
		priority   float64
		changeFreq sitemap.ChangeFreq
	}{
		{"/about", "About Us", 0.8, sitemap.Monthly},
		{"/products", "Our Products", 0.9, sitemap.Weekly},
		{"/blog", "Blog", 0.7, sitemap.Daily},
		{"/contact", "Contact Us", 0.5, sitemap.Yearly},
	}

	for _, page := range pages {
		sm.Add("http://localhost:8080"+page.url, time.Now(), page.priority, page.changeFreq,
			sitemap.WithTitle(page.title),
		)
	}

	// Add some blog posts with images
	blogPosts := []struct {
		slug  string
		title string
		image string
	}{
		{"first-post", "My First Blog Post", "/images/first-post.jpg"},
		{"second-post", "Another Great Post", "/images/second-post.jpg"},
		{"third-post", "Latest News", "/images/third-post.jpg"},
	}

	for _, post := range blogPosts {
		images := []sitemap.Image{
			{
				URL:   "http://localhost:8080" + post.image,
				Title: post.title + " Image",
			},
		}

		sm.Add("http://localhost:8080/blog/"+post.slug, time.Now(), 0.6, sitemap.Weekly,
			sitemap.WithTitle(post.title),
			sitemap.WithImages(images),
		)
	}

	return sm
}
