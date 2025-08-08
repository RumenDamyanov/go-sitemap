package fiberadapter

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.rumenx.com/sitemap"
)

func TestSitemap(t *testing.T) {
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
			checkBody:  func(body string) bool { return len(body) == 0 || body == "Internal Server Error" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fiber app
			app := fiber.New()
			app.Get("/sitemap.xml", Sitemap(tt.generator))

			// Create request
			req, err := http.NewRequest("GET", "/sitemap.xml", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Test request
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			// Check content type
			if tt.wantType != "" {
				contentType := resp.Header.Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			// Read and check body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if !tt.checkBody(string(body)) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapTXT(t *testing.T) {
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
			checkBody:  func(body string) bool { return len(body) == 0 || body == "Internal Server Error" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/sitemap.txt", SitemapTXT(tt.generator))

			req, err := http.NewRequest("GET", "/sitemap.txt", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantType != "" {
				contentType := resp.Header.Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if !tt.checkBody(string(body)) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapHTML(t *testing.T) {
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
			checkBody:  func(body string) bool { return len(body) == 0 || body == "Internal Server Error" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/sitemap.html", SitemapHTML(tt.generator))

			req, err := http.NewRequest("GET", "/sitemap.html", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantType != "" {
				contentType := resp.Header.Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if !tt.checkBody(string(body)) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}

func TestSitemapIndex(t *testing.T) {
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
			checkBody:  func(body string) bool { return len(body) == 0 || body == "Internal Server Error" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/sitemapindex.xml", SitemapIndex(tt.generator))

			req, err := http.NewRequest("GET", "/sitemapindex.xml", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantType != "" {
				contentType := resp.Header.Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if !tt.checkBody(string(body)) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}
