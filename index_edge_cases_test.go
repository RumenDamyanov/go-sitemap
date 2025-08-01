package sitemap

import (
	"strings"
	"testing"
	"time"
)

// Additional tests for uncovered code paths in index functionality

func TestIndexEdgeCases(t *testing.T) {
	idx := NewIndex()
	now := time.Now()

	// Test with invalid URLs
	err := idx.Add("invalid-url", now)
	if err == nil {
		t.Error("Index Add() should fail with invalid URL")
	}

	err = idx.Add("", now)
	if err == nil {
		t.Error("Index Add() should fail with empty URL")
	}

	err = idx.Add("ftp://example.com/sitemap.xml", now)
	if err == nil {
		t.Error("Index Add() should fail with invalid scheme")
	}

	// Test with valid URLs
	err = idx.Add("https://example.com/sitemap1.xml", now)
	if err != nil {
		t.Fatalf("Index Add() failed: %v", err)
	}

	if idx.Count() != 1 {
		t.Errorf("Expected 1 sitemap, got %d", idx.Count())
	}
}

func TestIndexXMLGeneration(t *testing.T) {
	idx := NewIndex()
	now := time.Now().Truncate(time.Second)

	// Test empty index
	xml, err := idx.XML()
	if err != nil {
		t.Fatalf("Empty index XML() failed: %v", err)
	}

	xmlStr := string(xml)
	if !strings.Contains(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Index XML should contain XML declaration")
	}

	if !strings.Contains(xmlStr, "sitemapindex") {
		t.Error("Index XML should contain sitemapindex element")
	}

	// Test with sitemaps
	err = idx.Add("https://example.com/sitemap1.xml", now)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	err = idx.Add("https://example.com/sitemap2.xml", now.AddDate(0, 0, -1))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	xml, err = idx.XML()
	if err != nil {
		t.Fatalf("Index XML() failed: %v", err)
	}

	xmlStr = string(xml)

	// Should contain both sitemaps
	if !strings.Contains(xmlStr, "sitemap1.xml") {
		t.Error("Index XML should contain first sitemap")
	}

	if !strings.Contains(xmlStr, "sitemap2.xml") {
		t.Error("Index XML should contain second sitemap")
	}

	// Should contain lastmod elements
	if !strings.Contains(xmlStr, "<lastmod>") {
		t.Error("Index XML should contain lastmod elements")
	}
}

func TestIndexWithZeroTime(t *testing.T) {
	idx := NewIndex()
	var zeroTime time.Time

	err := idx.Add("https://example.com/sitemap.xml", zeroTime)
	if err != nil {
		t.Fatalf("Add() with zero time failed: %v", err)
	}

	xml, err := idx.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)

	// Should not contain lastmod for zero time
	if strings.Contains(xmlStr, "<lastmod>") {
		t.Error("Index XML should not contain lastmod for zero time")
	}
}

func TestSitemapItemsMethod(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test empty sitemap
	items := sm.Items()
	if len(items) != 0 {
		t.Errorf("Empty sitemap should have 0 items, got %d", len(items))
	}

	// Add items and test
	sm.Add("https://example.com/1", now, 1.0, Daily)
	sm.Add("https://example.com/2", now, 0.8, Weekly)

	items = sm.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Note: Items() returns the actual slice (not a copy), so modifying it
	// would affect the original sitemap. This is current behavior.
	originalURL := items[0].URL
	if originalURL != "https://example.com/1" {
		t.Errorf("Expected first URL to be https://example.com/1, got %s", originalURL)
	}
}

func TestComplexImageAndVideoData(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with complex image data
	images := []Image{
		{
			URL:     "https://example.com/image1.jpg",
			Title:   "Image with special chars: <>\"&'",
			Caption: "Caption with\nnewlines and\ttabs",
		},
		{
			URL:   "https://example.com/image2.jpg",
			Title: "", // Empty title
		},
	}

	// Test with complex video data
	videos := []Video{
		{
			ThumbnailURL: "https://example.com/thumb1.jpg",
			Title:        "Video with special chars: <>\"&'",
			Description:  "Description with\nnewlines and\ttabs",
			ContentURL:   "https://example.com/video1.mp4",
			Duration:     3600, // 1 hour
		},
		{
			ThumbnailURL: "https://example.com/thumb2.jpg",
			Title:        "Video without content URL",
			Description:  "Simple description",
			// No ContentURL or Duration
		},
	}

	err := sm.Add(
		"https://example.com/complex",
		now,
		1.0,
		Daily,
		WithImages(images),
		WithVideos(videos),
	)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Test XML generation with complex data
	xml, err := sm.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)

	// XML package automatically escapes special characters when encoding
	// Check for properly encoded content
	if !strings.Contains(xmlStr, "image1.jpg") {
		t.Error("XML should contain image URLs")
	}

	// The encoding/xml package handles escaping automatically
	// Just verify the data is present
	if !strings.Contains(xmlStr, "Video with") {
		t.Error("XML should contain video titles")
	}

	// Test HTML generation with complex data
	html, err := sm.HTML()
	if err != nil {
		t.Fatalf("HTML() failed: %v", err)
	}

	htmlStr := string(html)

	// Should contain image and video data
	if !strings.Contains(htmlStr, "image1.jpg") {
		t.Error("HTML should contain image URLs")
	}

	if !strings.Contains(htmlStr, "Duration: 3600 seconds") {
		t.Error("HTML should contain video duration")
	}
}

func TestNewsWithComplexData(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with complex news data
	news := GoogleNews{
		SiteName:        "Site with special chars: <>\"&'",
		Language:        "en-US",
		PublicationDate: now,
		Title:           "Article with\nspecial\tchars: <>\"&'",
		Keywords:        "keyword1, keyword2, special chars: <>\"&'",
	}

	err := sm.Add("https://example.com/news", now, 1.0, Daily, WithGoogleNews(news))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Test GoogleNews format
	newsXML, err := sm.GoogleNews()
	if err != nil {
		t.Fatalf("GoogleNews() failed: %v", err)
	}

	xmlStr := string(newsXML)

	// XML encoding handles special characters automatically
	// Just verify the news data is present
	if !strings.Contains(xmlStr, "en-US") {
		t.Error("News XML should contain language")
	}

	if !strings.Contains(xmlStr, "Site with") {
		t.Error("News XML should contain site name")
	}
}
