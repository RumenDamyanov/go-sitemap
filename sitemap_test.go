package sitemap

import (
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	sm := New()
	if sm == nil {
		t.Fatal("New() returned nil")
	}

	if sm.Count() != 0 {
		t.Errorf("New sitemap should have 0 items, got %d", sm.Count())
	}

	if sm.opts.MaxURLs != 50000 {
		t.Errorf("Default MaxURLs should be 50000, got %d", sm.opts.MaxURLs)
	}
}

func TestNewWithOptions(t *testing.T) {
	opts := &Options{
		MaxURLs:     1000,
		BaseURL:     "https://example.com",
		PreAllocate: true,
	}

	sm := NewWithOptions(opts)
	if sm == nil {
		t.Fatal("NewWithOptions() returned nil")
	}

	if sm.opts.MaxURLs != 1000 {
		t.Errorf("MaxURLs should be 1000, got %d", sm.opts.MaxURLs)
	}

	if sm.opts.BaseURL != "https://example.com" {
		t.Errorf("BaseURL should be 'https://example.com', got %s", sm.opts.BaseURL)
	}
}

func TestAdd(t *testing.T) {
	sm := New()
	now := time.Now()

	err := sm.Add("https://example.com/", now, 1.0, Daily)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if sm.Count() != 1 {
		t.Errorf("Expected 1 item, got %d", sm.Count())
	}

	items := sm.Items()
	if len(items) != 1 {
		t.Errorf("Expected 1 item in Items(), got %d", len(items))
	}

	item := items[0]
	if item.URL != "https://example.com/" {
		t.Errorf("Expected URL 'https://example.com/', got %s", item.URL)
	}

	if item.Priority != 1.0 {
		t.Errorf("Expected priority 1.0, got %f", item.Priority)
	}

	if item.ChangeFreq != Daily {
		t.Errorf("Expected changefreq 'daily', got %s", item.ChangeFreq)
	}
}

func TestAddInvalidURL(t *testing.T) {
	sm := New()
	now := time.Now()

	tests := []struct {
		name string
		url  string
	}{
		{"empty URL", ""},
		{"relative URL", "/relative"},
		{"invalid scheme", "ftp://example.com"},
		{"malformed URL", "http://[::1:80"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.Add(tt.url, now, 1.0, Daily)
			if err == nil {
				t.Errorf("Add() should have failed for %s", tt.url)
			}
		})
	}
}

func TestAddInvalidPriority(t *testing.T) {
	sm := New()
	now := time.Now()

	tests := []struct {
		priority float64
	}{
		{-0.1},
		{1.1},
		{2.0},
		{-1.0},
	}

	for _, tt := range tests {
		err := sm.Add("https://example.com/", now, tt.priority, Daily)
		if err == nil {
			t.Errorf("Add() should have failed for priority %f", tt.priority)
		}
	}
}

func TestAddWithOptions(t *testing.T) {
	sm := New()
	now := time.Now()

	images := []Image{
		{URL: "https://example.com/image.jpg", Title: "Test Image"},
	}

	videos := []Video{
		{ThumbnailURL: "https://example.com/thumb.jpg", Title: "Test Video", Description: "A test video"},
	}

	news := GoogleNews{
		SiteName:        "Test News",
		Language:        "en",
		PublicationDate: now,
		Title:           "Test Article",
	}

	err := sm.Add(
		"https://example.com/article",
		now,
		0.8,
		Weekly,
		WithTitle("Test Article"),
		WithImages(images),
		WithVideos(videos),
		WithGoogleNews(news),
	)

	if err != nil {
		t.Fatalf("Add() with options failed: %v", err)
	}

	items := sm.Items()
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %s", item.Title)
	}

	if len(item.Images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(item.Images))
	}

	if len(item.Videos) != 1 {
		t.Errorf("Expected 1 video, got %d", len(item.Videos))
	}

	if item.News == nil {
		t.Error("Expected news metadata, got nil")
	} else if item.News.SiteName != "Test News" {
		t.Errorf("Expected news site name 'Test News', got %s", item.News.SiteName)
	}
}

func TestAddItem(t *testing.T) {
	sm := New()
	now := time.Now()

	item := Item{
		URL:        "https://example.com/page",
		LastMod:    now,
		Priority:   0.7,
		ChangeFreq: Monthly,
		Title:      "Test Page",
	}

	err := sm.AddItem(item)
	if err != nil {
		t.Fatalf("AddItem() failed: %v", err)
	}

	if sm.Count() != 1 {
		t.Errorf("Expected 1 item, got %d", sm.Count())
	}

	items := sm.Items()
	addedItem := items[0]
	if addedItem.URL != item.URL {
		t.Errorf("Expected URL %s, got %s", item.URL, addedItem.URL)
	}

	if addedItem.Title != item.Title {
		t.Errorf("Expected title %s, got %s", item.Title, addedItem.Title)
	}
}

func TestAddItems(t *testing.T) {
	sm := New()

	items := []Item{
		{URL: "https://example.com/page1", Priority: 0.8, ChangeFreq: Weekly},
		{URL: "https://example.com/page2", Priority: 0.7, ChangeFreq: Monthly},
		{URL: "https://example.com/page3", Priority: 0.6, ChangeFreq: Yearly},
	}

	err := sm.AddItems(items)
	if err != nil {
		t.Fatalf("AddItems() failed: %v", err)
	}

	if sm.Count() != 3 {
		t.Errorf("Expected 3 items, got %d", sm.Count())
	}

	sitemapItems := sm.Items()
	for i, expected := range items {
		if sitemapItems[i].URL != expected.URL {
			t.Errorf("Item %d: expected URL %s, got %s", i, expected.URL, sitemapItems[i].URL)
		}
	}
}

func TestMaxURLsLimit(t *testing.T) {
	opts := &Options{MaxURLs: 2}
	sm := NewWithOptions(opts)
	now := time.Now()

	// Add max items
	for i := 0; i < 2; i++ {
		err := sm.Add("https://example.com/", now, 1.0, Daily)
		if err != nil {
			t.Fatalf("Add() failed for item %d: %v", i, err)
		}
	}

	// This should fail
	err := sm.Add("https://example.com/", now, 1.0, Daily)
	if err == nil {
		t.Error("Add() should have failed when exceeding MaxURLs limit")
	}
}

func TestClear(t *testing.T) {
	sm := New()
	now := time.Now()

	// Add some items
	sm.Add("https://example.com/1", now, 1.0, Daily)
	sm.Add("https://example.com/2", now, 1.0, Daily)

	if sm.Count() != 2 {
		t.Errorf("Expected 2 items before clear, got %d", sm.Count())
	}

	sm.Clear()

	if sm.Count() != 0 {
		t.Errorf("Expected 0 items after clear, got %d", sm.Count())
	}
}

func TestXML(t *testing.T) {
	sm := New()
	now := time.Now().Truncate(time.Second) // Remove nanoseconds for comparison

	err := sm.Add("https://example.com/", now, 1.0, Daily)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	xml, err := sm.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)

	// Check XML declaration
	if !strings.Contains(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("XML should contain XML declaration")
	}

	// Check namespace
	if !strings.Contains(xmlStr, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("XML should contain sitemap namespace")
	}

	// Check URL
	if !strings.Contains(xmlStr, "https://example.com/") {
		t.Error("XML should contain the URL")
	}

	// Check priority
	if !strings.Contains(xmlStr, "<priority>1.0</priority>") {
		t.Error("XML should contain priority")
	}

	// Check changefreq
	if !strings.Contains(xmlStr, "<changefreq>daily</changefreq>") {
		t.Error("XML should contain changefreq")
	}
}

func TestTXT(t *testing.T) {
	sm := New()
	now := time.Now()

	urls := []string{
		"https://example.com/",
		"https://example.com/about",
		"https://example.com/contact",
	}

	for _, url := range urls {
		err := sm.Add(url, now, 1.0, Daily)
		if err != nil {
			t.Fatalf("Add() failed for %s: %v", url, err)
		}
	}

	txt, err := sm.TXT()
	if err != nil {
		t.Fatalf("TXT() failed: %v", err)
	}

	txtStr := string(txt)
	for _, url := range urls {
		if !strings.Contains(txtStr, url) {
			t.Errorf("TXT should contain URL %s", url)
		}
	}

	lines := strings.Split(strings.TrimSpace(txtStr), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines in TXT output, got %d", len(lines))
	}
}

func TestHTML(t *testing.T) {
	sm := New()
	now := time.Now()

	err := sm.Add("https://example.com/", now, 1.0, Daily, WithTitle("Homepage"))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	html, err := sm.HTML()
	if err != nil {
		t.Fatalf("HTML() failed: %v", err)
	}

	htmlStr := string(html)

	// Check HTML structure
	if !strings.Contains(htmlStr, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE")
	}

	if !strings.Contains(htmlStr, "<title>Sitemap</title>") {
		t.Error("HTML should contain title")
	}

	if !strings.Contains(htmlStr, "https://example.com/") {
		t.Error("HTML should contain the URL")
	}

	if !strings.Contains(htmlStr, "Homepage") {
		t.Error("HTML should contain the page title")
	}
}

func TestChangeFreqConstants(t *testing.T) {
	tests := []struct {
		freq     ChangeFreq
		expected string
	}{
		{Always, "always"},
		{Hourly, "hourly"},
		{Daily, "daily"},
		{Weekly, "weekly"},
		{Monthly, "monthly"},
		{Yearly, "yearly"},
		{Never, "never"},
	}

	for _, tt := range tests {
		if string(tt.freq) != tt.expected {
			t.Errorf("ChangeFreq %v should be %s, got %s", tt.freq, tt.expected, string(tt.freq))
		}
	}
}
