package sitemap

import (
	"strings"
	"testing"
	"time"
)

// Test edge cases and error conditions that might not be covered

func TestAddWithInvalidOptions(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with invalid image URL (should still add the item but with validation)
	images := []Image{
		{URL: "not-a-valid-url", Title: "Invalid Image"},
	}

	err := sm.Add("https://example.com/", now, 1.0, Daily, WithImages(images))
	if err != nil {
		t.Fatalf("Add() should not fail with invalid image URL: %v", err)
	}

	// Item should still be added
	if sm.Count() != 1 {
		t.Errorf("Expected 1 item, got %d", sm.Count())
	}
}

func TestAddWithEmptyOptions(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with empty slices
	err := sm.Add(
		"https://example.com/",
		now,
		1.0,
		Daily,
		WithImages([]Image{}),
		WithVideos([]Video{}),
	)
	if err != nil {
		t.Fatalf("Add() failed with empty options: %v", err)
	}

	items := sm.Items()
	item := items[0]
	if len(item.Images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(item.Images))
	}
	if len(item.Videos) != 0 {
		t.Errorf("Expected 0 videos, got %d", len(item.Videos))
	}
}

func TestURLValidation(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test various URL validation scenarios
	tests := []struct {
		name      string
		url       string
		shouldErr bool
	}{
		{"valid HTTPS", "https://example.com/path", false},
		{"valid HTTP", "http://example.com/path", false},
		{"with query params", "https://example.com/path?param=value", false},
		{"with fragment", "https://example.com/path#section", false},
		{"with port", "https://example.com:8080/path", false},
		{"invalid scheme", "ftp://example.com/", true},
		{"no scheme", "example.com/path", true},
		{"empty", "", true},
		{"malformed", "http://[::1:80", true},
		{"relative path", "/relative/path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm.Clear() // Clear for each test
			err := sm.Add(tt.url, now, 1.0, Daily)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for URL %s", tt.url)
			} else if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for URL %s: %v", tt.url, err)
			}
		})
	}
}

func TestPriorityValidation(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test boundary conditions
	tests := []struct {
		priority  float64
		shouldErr bool
	}{
		{0.0, false}, // minimum valid
		{1.0, false}, // maximum valid
		{0.5, false}, // middle value
		{-0.1, true}, // below minimum
		{1.1, true},  // above maximum
		{-1.0, true}, // negative
		{2.0, true},  // too high
	}

	for _, tt := range tests {
		sm.Clear()
		err := sm.Add("https://example.com/", now, tt.priority, Daily)
		if tt.shouldErr && err == nil {
			t.Errorf("Expected error for priority %f", tt.priority)
		} else if !tt.shouldErr && err != nil {
			t.Errorf("Unexpected error for priority %f: %v", tt.priority, err)
		}
	}
}

func TestAddItemValidation(t *testing.T) {
	sm := New()

	// Test AddItem with invalid item
	invalidItem := Item{
		URL:      "invalid-url",
		Priority: 2.0, // Invalid priority
	}

	err := sm.AddItem(invalidItem)
	if err == nil {
		t.Error("AddItem should fail with invalid item")
	}

	// Test AddItem with valid item
	validItem := Item{
		URL:        "https://example.com/",
		Priority:   0.8,
		ChangeFreq: Weekly,
		Title:      "Test Page",
	}

	err = sm.AddItem(validItem)
	if err != nil {
		t.Fatalf("AddItem should succeed with valid item: %v", err)
	}
}

func TestAddItemsValidation(t *testing.T) {
	sm := New()

	// Test AddItems with mix of valid and invalid items
	items := []Item{
		{URL: "https://example.com/1", Priority: 0.8, ChangeFreq: Weekly},
		{URL: "invalid-url", Priority: 0.8, ChangeFreq: Weekly}, // Invalid URL
		{URL: "https://example.com/3", Priority: 0.8, ChangeFreq: Weekly},
	}

	err := sm.AddItems(items)
	if err == nil {
		t.Error("AddItems should fail with invalid item in the list")
	}

	// Should have added the first valid item before failing on the second
	if sm.Count() != 1 {
		t.Errorf("Expected 1 item after failed AddItems (first item should be added), got %d", sm.Count())
	}

	// Clear and test with all valid items
	sm.Clear()
	validItems := []Item{
		{URL: "https://example.com/1", Priority: 0.8, ChangeFreq: Weekly},
		{URL: "https://example.com/2", Priority: 0.7, ChangeFreq: Monthly},
		{URL: "https://example.com/3", Priority: 0.6, ChangeFreq: Yearly},
	}

	err = sm.AddItems(validItems)
	if err != nil {
		t.Fatalf("AddItems should succeed with valid items: %v", err)
	}

	if sm.Count() != 3 {
		t.Errorf("Expected 3 items, got %d", sm.Count())
	}
}

func TestOptionsWithNil(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with nil GoogleNews
	err := sm.Add(
		"https://example.com/",
		now,
		1.0,
		Daily,
		WithGoogleNews(GoogleNews{}), // Empty but not nil
	)
	if err != nil {
		t.Fatalf("Add() failed with empty GoogleNews: %v", err)
	}

	items := sm.Items()
	if items[0].News == nil {
		t.Error("Expected non-nil news, got nil")
	}
}

func TestLargeDatasets(t *testing.T) {
	sm := New()
	now := time.Now()

	// Test with many items to ensure performance is acceptable
	const numItems = 1000

	for i := 0; i < numItems; i++ {
		url := "https://example.com/page" + string(rune('0'+i%10))
		err := sm.Add(url, now, 0.5, Daily)
		if err != nil {
			t.Fatalf("Add() failed for item %d: %v", i, err)
		}
	}

	if sm.Count() != numItems {
		t.Errorf("Expected %d items, got %d", numItems, sm.Count())
	}

	// Test XML generation with large dataset
	xml, err := sm.XML()
	if err != nil {
		t.Fatalf("XML() failed with large dataset: %v", err)
	}

	if len(xml) == 0 {
		t.Error("XML should not be empty")
	}

	// Verify XML contains expected number of URLs
	xmlStr := string(xml)
	urlCount := strings.Count(xmlStr, "<loc>")
	if urlCount != numItems {
		t.Errorf("Expected %d URLs in XML, found %d", numItems, urlCount)
	}
}

func TestZeroTime(t *testing.T) {
	sm := New()

	// Test with zero time
	var zeroTime time.Time
	err := sm.Add("https://example.com/", zeroTime, 1.0, Daily)
	if err != nil {
		t.Fatalf("Add() failed with zero time: %v", err)
	}

	items := sm.Items()
	item := items[0]
	if !item.LastMod.IsZero() {
		t.Error("Expected zero time to remain zero")
	}

	// Test XML output with zero time
	xml, err := sm.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)
	if strings.Contains(xmlStr, "<lastmod>") {
		t.Error("XML should not contain lastmod element for zero time")
	}
}

func TestChangeFreqValues(t *testing.T) {
	sm := New()
	now := time.Now()

	freqs := []ChangeFreq{Always, Hourly, Daily, Weekly, Monthly, Yearly, Never}

	for i, freq := range freqs {
		url := "https://example.com/page" + string(rune('0'+i))
		err := sm.Add(url, now, 1.0, freq)
		if err != nil {
			t.Fatalf("Add() failed for freq %s: %v", freq, err)
		}
	}

	xml, err := sm.XML()
	if err != nil {
		t.Fatalf("XML() failed: %v", err)
	}

	xmlStr := string(xml)
	for _, freq := range freqs {
		if !strings.Contains(xmlStr, "<changefreq>"+string(freq)+"</changefreq>") {
			t.Errorf("XML should contain changefreq %s", freq)
		}
	}
}

func TestBaseURLWithOptions(t *testing.T) {
	opts := &Options{
		BaseURL: "https://example.com",
	}
	sm := NewWithOptions(opts)
	now := time.Now()

	// Test relative URL that should work with base URL
	err := sm.Add("/relative/path", now, 1.0, Daily)
	if err == nil {
		t.Error("Add() should still validate URLs even with BaseURL option")
	}

	// Base URL option doesn't automatically make relative URLs valid
	// It's more for internal tracking/validation
}

func TestPreAllocateOption(t *testing.T) {
	opts := &Options{
		MaxURLs:     100,
		PreAllocate: true,
	}
	sm := NewWithOptions(opts)

	// The PreAllocate option should pre-allocate the slice
	// This is mainly for performance, hard to test directly
	if sm == nil {
		t.Error("NewWithOptions should not return nil")
	}

	// Add some items to test it works normally
	now := time.Now()
	for i := 0; i < 10; i++ {
		url := "https://example.com/page" + string(rune('0'+i))
		err := sm.Add(url, now, 1.0, Daily)
		if err != nil {
			t.Fatalf("Add() failed: %v", err)
		}
	}

	if sm.Count() != 10 {
		t.Errorf("Expected 10 items, got %d", sm.Count())
	}
}
