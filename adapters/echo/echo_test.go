package echoadapter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-sitemap"
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
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Echo instance
			e := echo.New()
			e.GET("/sitemap.xml", Sitemap(tt.generator))

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
			rec := httptest.NewRecorder()

			// Serve request
			e.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			// Check content type
			if tt.wantType != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			// Check body
			if !tt.checkBody(rec.Body.String()) {
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
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/sitemap.txt", SitemapTXT(tt.generator))

			req := httptest.NewRequest(http.MethodGet, "/sitemap.txt", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantType != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(rec.Body.String()) {
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
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/sitemap.html", SitemapHTML(tt.generator))

			req := httptest.NewRequest(http.MethodGet, "/sitemap.html", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantType != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(rec.Body.String()) {
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
			checkBody:  func(body string) bool { return body == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/sitemapindex.xml", SitemapIndex(tt.generator))

			req := httptest.NewRequest(http.MethodGet, "/sitemapindex.xml", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantType != "" {
				contentType := rec.Header().Get("Content-Type")
				if contentType != tt.wantType {
					t.Errorf("Expected content type %s, got %s", tt.wantType, contentType)
				}
			}

			if !tt.checkBody(rec.Body.String()) {
				t.Errorf("Body check failed for test %s", tt.name)
			}
		})
	}
}
