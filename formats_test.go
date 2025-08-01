package sitemap

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestJSON(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test empty sitemap
	jsonData, err := sm.JSON()
	if err != nil {
		t.Fatalf("JSON() failed: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["count"] != float64(0) {
		t.Errorf("Expected count 0, got %v", result["count"])
	}

	// Test with items
	err = sm.Add("https://example.com/", now, 1.0, Daily, WithTitle("Homepage"))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	err = sm.Add("https://example.com/about", now, 0.8, Weekly)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	jsonData, err = sm.JSON()
	if err != nil {
		t.Fatalf("JSON() failed: %v", err)
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["count"] != float64(2) {
		t.Errorf("Expected count 2, got %v", result["count"])
	}

	urls, ok := result["urls"].([]interface{})
	if !ok {
		t.Fatal("URLs field should be an array")
	}

	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(urls))
	}

	// Verify JSON is properly formatted
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "https://example.com/") {
		t.Error("JSON should contain the first URL")
	}

	if !strings.Contains(jsonStr, "https://example.com/about") {
		t.Error("JSON should contain the second URL")
	}
}

func TestGoogleNews(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with no news items
	newsXML, err := sm.GoogleNews()
	if err != nil {
		t.Fatalf("GoogleNews() failed: %v", err)
	}

	// Should contain empty sitemap structure
	xmlStr := string(newsXML)
	if !strings.Contains(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("News XML should contain XML declaration")
	}

	// Test with news items
	news1 := GoogleNews{
		SiteName:        "Test News Site",
		Language:        "en",
		PublicationDate: now,
		Title:           "Breaking News Article",
		Keywords:        "news, breaking, test",
	}

	news2 := GoogleNews{
		SiteName:        "Another News Site",
		Language:        "es",
		PublicationDate: now.AddDate(0, 0, -1),
		Title:           "Another Article",
	}

	// Add regular item (should not appear in news sitemap)
	err = sm.Add("https://example.com/regular", now, 1.0, Daily)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Add news items
	err = sm.Add("https://example.com/news1", now, 1.0, Daily, WithGoogleNews(news1))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	err = sm.Add("https://example.com/news2", now, 0.9, Hourly, WithGoogleNews(news2))
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	newsXML, err = sm.GoogleNews()
	if err != nil {
		t.Fatalf("GoogleNews() failed: %v", err)
	}

	xmlStr = string(newsXML)

	// Should contain only news items
	if !strings.Contains(xmlStr, "https://example.com/news1") {
		t.Error("News XML should contain first news URL")
	}

	if !strings.Contains(xmlStr, "https://example.com/news2") {
		t.Error("News XML should contain second news URL")
	}

	// Should NOT contain regular URL
	if strings.Contains(xmlStr, "https://example.com/regular") {
		t.Error("News XML should not contain regular (non-news) URL")
	}

	// Should contain news metadata
	if !strings.Contains(xmlStr, "Breaking News Article") {
		t.Error("News XML should contain news title")
	}

	if !strings.Contains(xmlStr, "Test News Site") {
		t.Error("News XML should contain news site name")
	}

	if !strings.Contains(xmlStr, "news, breaking, test") {
		t.Error("News XML should contain keywords")
	}
}

func TestMobile(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test empty sitemap
	mobileXML, err := sm.Mobile()
	if err != nil {
		t.Fatalf("Mobile() failed: %v", err)
	}

	xmlStr := string(mobileXML)
	if !strings.Contains(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Mobile XML should contain XML declaration")
	}

	if !strings.Contains(xmlStr, `xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0"`) {
		t.Error("Mobile XML should contain mobile namespace")
	}

	// Test with items
	err = sm.Add("https://m.example.com/", now, 1.0, Daily)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	err = sm.Add("https://m.example.com/products", now, 0.8, Weekly)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	mobileXML, err = sm.Mobile()
	if err != nil {
		t.Fatalf("Mobile() failed: %v", err)
	}

	xmlStr = string(mobileXML)

	// Should contain URLs
	if !strings.Contains(xmlStr, "https://m.example.com/") {
		t.Error("Mobile XML should contain first URL")
	}

	if !strings.Contains(xmlStr, "https://m.example.com/products") {
		t.Error("Mobile XML should contain second URL")
	}

	// Test XML structure - check that namespaces are present in the XML string
	if !strings.Contains(xmlStr, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("Mobile XML should contain standard sitemap namespace")
	}

	if !strings.Contains(xmlStr, `xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0"`) {
		t.Error("Mobile XML should contain mobile namespace")
	}

	// Count URLs by parsing XML structure
	var urlset struct {
		XMLName xml.Name `xml:"urlset"`
		URLs    []struct {
			URL string `xml:"loc"`
		} `xml:"url"`
	}

	err = xml.Unmarshal(mobileXML, &urlset)
	if err != nil {
		t.Fatalf("Failed to unmarshal mobile XML: %v", err)
	}

	if len(urlset.URLs) != 2 {
		t.Errorf("Expected 2 URLs in mobile sitemap, got %d", len(urlset.URLs))
	}
}

func TestHTMLWithComplexData(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with images, videos, and news
	images := []Image{
		{URL: "https://example.com/image1.jpg", Title: "Test Image 1", Caption: "A test image"},
		{URL: "https://example.com/image2.jpg", Title: "Test Image 2"},
	}

	videos := []Video{
		{
			ThumbnailURL: "https://example.com/thumb1.jpg",
			Title:        "Test Video 1",
			Description:  "A test video description",
			ContentURL:   "https://example.com/video1.mp4",
			Duration:     120,
		},
		{
			ThumbnailURL: "https://example.com/thumb2.jpg",
			Title:        "Test Video 2",
			Description:  "Another test video",
		},
	}

	news := GoogleNews{
		SiteName:        "Test News",
		Language:        "en",
		PublicationDate: now,
		Title:           "Test Article",
		Keywords:        "test, article",
	}

	err := sm.Add(
		"https://example.com/complex",
		now,
		0.9,
		Weekly,
		WithTitle("Complex Page"),
		WithImages(images),
		WithVideos(videos),
		WithGoogleNews(news),
	)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	html, err := sm.HTML()
	if err != nil {
		t.Fatalf("HTML() failed: %v", err)
	}

	htmlStr := string(html)

	// Check for complex data in HTML
	if !strings.Contains(htmlStr, "Complex Page") {
		t.Error("HTML should contain page title")
	}

	if !strings.Contains(htmlStr, "Test Image 1") {
		t.Error("HTML should contain first image title")
	}

	if !strings.Contains(htmlStr, "A test image") {
		t.Error("HTML should contain image caption")
	}

	if !strings.Contains(htmlStr, "Test Video 1") {
		t.Error("HTML should contain first video title")
	}

	if !strings.Contains(htmlStr, "Duration: 120 seconds") {
		t.Error("HTML should contain video duration")
	}

	if !strings.Contains(htmlStr, "Test News") {
		t.Error("HTML should contain news site name")
	}

	if !strings.Contains(htmlStr, "test, article") {
		t.Error("HTML should contain news keywords")
	}

	// Check stats section
	if !strings.Contains(htmlStr, "Total URLs:</strong> 1") {
		t.Error("HTML should show correct URL count")
	}
}

func TestHTMLTemplateError(t *testing.T) {
	// This test verifies that the HTML method handles template parsing errors gracefully
	// We can't easily trigger a template error with the current implementation,
	// but we test the basic functionality
	sm := New()

	html, err := sm.HTML()
	if err != nil {
		t.Fatalf("HTML() failed: %v", err)
	}

	htmlStr := string(html)
	if !strings.Contains(htmlStr, "Total URLs:</strong> 0") {
		t.Error("Empty sitemap HTML should show 0 URLs")
	}
}

func TestTXTWithManyURLs(t *testing.T) {
	sm := New()
	now := time.Now()

	urls := make([]string, 100)
	for i := 0; i < 100; i++ {
		urls[i] = "https://example.com/page" + string(rune('0'+i%10))
		err := sm.Add(urls[i], now, 1.0, Daily)
		if err != nil {
			t.Fatalf("Add() failed for URL %d: %v", i, err)
		}
	}

	txt, err := sm.TXT()
	if err != nil {
		t.Fatalf("TXT() failed: %v", err)
	}

	txtStr := string(txt)
	lines := strings.Split(strings.TrimSpace(txtStr), "\n")

	if len(lines) != 100 {
		t.Errorf("Expected 100 lines in TXT output, got %d", len(lines))
	}

	// Verify each line is a valid URL
	for i, line := range lines {
		if !strings.HasPrefix(line, "https://example.com/page") {
			t.Errorf("Line %d should start with https://example.com/page, got %s", i, line)
		}
	}
}
