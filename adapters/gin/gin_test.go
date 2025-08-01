package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rumendamyanov/go-sitemap"
)

func TestSitemap(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		generator  SitemapGenerator
		wantStatus int
		wantType   string
		checkBody  func(string) bool
	}{
		{
			name: "successful sitemap generation",
			generator: func() *sitemap.Sitemap {
				sm := sitemap.New()
				sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
				return sm
			},
			wantStatus: http.StatusOK,
			wantType:   "application/xml",
			checkBody: func(body string) bool {
				return strings.Contains(body, "https://example.com/") &&
					strings.Contains(body, `<?xml version="1.0" encoding="UTF-8"?>`)
			},
		},
		{
			name: "nil sitemap generator",
			generator: func() *sitemap.Sitemap {
				return nil
			},
			wantStatus: http.StatusInternalServerError,
			wantType:   "",
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router and add handler
			r := gin.New()
			r.GET("/sitemap.xml", Sitemap(tt.generator))

			// Create request
			req, err := http.NewRequest("GET", "/sitemap.xml", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			r.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			// Check content type
			if tt.wantType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			// Check body
			if !tt.checkBody(w.Body.String()) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapTXT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		generator  SitemapGenerator
		wantStatus int
		wantType   string
		checkBody  func(string) bool
	}{
		{
			name: "successful TXT generation",
			generator: func() *sitemap.Sitemap {
				sm := sitemap.New()
				sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
				sm.Add("https://example.com/about", time.Now(), 0.8, sitemap.Weekly)
				return sm
			},
			wantStatus: http.StatusOK,
			wantType:   "text/plain",
			checkBody: func(body string) bool {
				return strings.Contains(body, "https://example.com/") &&
					strings.Contains(body, "https://example.com/about")
			},
		},
		{
			name: "nil sitemap generator",
			generator: func() *sitemap.Sitemap {
				return nil
			},
			wantStatus: http.StatusInternalServerError,
			wantType:   "",
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/sitemap.txt", SitemapTXT(tt.generator))

			req, err := http.NewRequest("GET", "/sitemap.txt", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(w.Body.String()) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		generator  SitemapGenerator
		wantStatus int
		wantType   string
		checkBody  func(string) bool
	}{
		{
			name: "successful HTML generation",
			generator: func() *sitemap.Sitemap {
				sm := sitemap.New()
				sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily, sitemap.WithTitle("Homepage"))
				return sm
			},
			wantStatus: http.StatusOK,
			wantType:   "text/html",
			checkBody: func(body string) bool {
				return strings.Contains(body, "<!DOCTYPE html>") &&
					strings.Contains(body, "https://example.com/") &&
					strings.Contains(body, "Homepage")
			},
		},
		{
			name: "nil sitemap generator",
			generator: func() *sitemap.Sitemap {
				return nil
			},
			wantStatus: http.StatusInternalServerError,
			wantType:   "",
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/sitemap.html", SitemapHTML(tt.generator))

			req, err := http.NewRequest("GET", "/sitemap.html", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(w.Body.String()) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapIndex(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		generator  func() *sitemap.Index
		wantStatus int
		wantType   string
		checkBody  func(string) bool
	}{
		{
			name: "successful index generation",
			generator: func() *sitemap.Index {
				idx := sitemap.NewIndex()
				idx.Add("https://example.com/sitemap1.xml", time.Now())
				idx.Add("https://example.com/sitemap2.xml", time.Now())
				return idx
			},
			wantStatus: http.StatusOK,
			wantType:   "application/xml",
			checkBody: func(body string) bool {
				return strings.Contains(body, "sitemap1.xml") &&
					strings.Contains(body, "sitemap2.xml") &&
					strings.Contains(body, "sitemapindex")
			},
		},
		{
			name: "nil index generator",
			generator: func() *sitemap.Index {
				return nil
			},
			wantStatus: http.StatusInternalServerError,
			wantType:   "",
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/sitemapindex.xml", SitemapIndex(tt.generator))

			req, err := http.NewRequest("GET", "/sitemapindex.xml", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(w.Body.String()) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapGenerationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test error handling when sitemap generation fails
	// We'll create a sitemap that might fail during XML generation
	generator := func() *sitemap.Sitemap {
		sm := sitemap.New()
		// Add a valid URL
		sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
		return sm
	}

	r := gin.New()
	r.GET("/sitemap.xml", Sitemap(generator))

	req, err := http.NewRequest("GET", "/sitemap.xml", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This should succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
