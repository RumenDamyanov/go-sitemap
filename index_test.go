package sitemap

import (
	"strings"
	"testing"
	"time"
)

func TestNewIndex(t *testing.T) {
	idx := NewIndex()
	if idx == nil {
		t.Fatal("NewIndex() returned nil")
	}

	if idx.Count() != 0 {
		t.Errorf("New index should have 0 sitemaps, got %d", idx.Count())
	}
}

func TestIndexAdd(t *testing.T) {
	idx := NewIndex()
	now := time.Now()

	err := idx.Add("https://example.com/sitemap1.xml", now)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if idx.Count() != 1 {
		t.Errorf("Expected 1 sitemap, got %d", idx.Count())
	}

	err = idx.Add("https://example.com/sitemap2.xml", now)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if idx.Count() != 2 {
		t.Errorf("Expected 2 sitemaps, got %d", idx.Count())
	}
}

func TestIndexAddInvalidURL(t *testing.T) {
	idx := NewIndex()
	now := time.Now()

	tests := []struct {
		name string
		url  string
	}{
		{"empty URL", ""},
		{"relative URL", "/sitemap.xml"},
		{"invalid scheme", "ftp://example.com/sitemap.xml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := idx.Add(tt.url, now)
			if err == nil {
				t.Errorf("Add() should have failed for %s", tt.url)
			}
		})
	}
}

func TestIndexXML(t *testing.T) {
	idx := NewIndex()
	now := time.Now().Truncate(time.Second)

	sitemaps := []string{
		"https://example.com/sitemap-products.xml",
		"https://example.com/sitemap-blog.xml",
		"https://example.com/sitemap-news.xml",
	}

	for _, sitemapURL := range sitemaps {
		err := idx.Add(sitemapURL, now)
		if err != nil {
			t.Fatalf("Add() failed for %s: %v", sitemapURL, err)
		}
	}

	xml, err := idx.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)

	// Check XML declaration
	if !strings.Contains(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("XML should contain XML declaration")
	}

	// Check root element
	if !strings.Contains(xmlStr, "<sitemapindex") {
		t.Error("XML should contain sitemapindex root element")
	}

	// Check namespace
	if !strings.Contains(xmlStr, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("XML should contain sitemap namespace")
	}

	// Check all sitemap URLs
	for _, sitemapURL := range sitemaps {
		if !strings.Contains(xmlStr, sitemapURL) {
			t.Errorf("XML should contain sitemap URL %s", sitemapURL)
		}
	}

	// Check lastmod is present
	expectedLastMod := now.Format(time.RFC3339)
	if !strings.Contains(xmlStr, expectedLastMod) {
		t.Errorf("XML should contain lastmod %s", expectedLastMod)
	}
}
