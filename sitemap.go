// Package sitemap provides functionality for generating XML sitemaps
// following the sitemaps.org protocol. It supports standard sitemaps,
// Google News sitemaps, image sitemaps, and video sitemaps.
//
// The package is designed to be framework-agnostic and includes
// adapters for popular Go web frameworks.
package sitemap

import (
	"fmt"
	"net/url"
	"time"
)

// ChangeFreq represents how frequently the page is likely to change.
type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

// Sitemap represents a sitemap that can contain multiple URLs with their metadata.
type Sitemap struct {
	items []Item
	opts  Options
}

// Options contains configuration options for the sitemap.
type Options struct {
	MaxURLs     int
	BaseURL     string
	PreAllocate bool
}

// Item represents a single URL entry in the sitemap.
type Item struct {
	URL        string        `xml:"loc" json:"url"`
	LastMod    time.Time     `xml:"lastmod,omitempty" json:"lastmod,omitempty"`
	ChangeFreq ChangeFreq    `xml:"changefreq,omitempty" json:"changefreq,omitempty"`
	Priority   float64       `xml:"priority,omitempty" json:"priority,omitempty"`
	Title      string        `xml:"-" json:"title,omitempty"`
	Images     []Image       `xml:"image:image,omitempty" json:"images,omitempty"`
	Videos     []Video       `xml:"video:video,omitempty" json:"videos,omitempty"`
	News       *GoogleNews   `xml:"news:news,omitempty" json:"news,omitempty"`
	Alternates []Alternate   `xml:"-" json:"alternates,omitempty"`
	Langs      []Translation `xml:"-" json:"translations,omitempty"`
}

// Image represents an image reference in a sitemap entry.
type Image struct {
	URL     string `xml:"image:loc" json:"url"`
	Title   string `xml:"image:title,omitempty" json:"title,omitempty"`
	Caption string `xml:"image:caption,omitempty" json:"caption,omitempty"`
}

// Video represents a video reference in a sitemap entry.
type Video struct {
	ThumbnailURL string `xml:"video:thumbnail_loc" json:"thumbnail_url"`
	Title        string `xml:"video:title" json:"title"`
	Description  string `xml:"video:description" json:"description"`
	ContentURL   string `xml:"video:content_loc,omitempty" json:"content_url,omitempty"`
	PlayerURL    string `xml:"video:player_loc,omitempty" json:"player_url,omitempty"`
	Duration     int    `xml:"video:duration,omitempty" json:"duration,omitempty"`
}

// GoogleNews represents Google News specific metadata.
type GoogleNews struct {
	SiteName        string    `xml:"news:publication>news:name" json:"site_name"`
	Language        string    `xml:"news:publication>news:language" json:"language"`
	PublicationDate time.Time `xml:"news:publication_date" json:"publication_date"`
	Title           string    `xml:"news:title" json:"title"`
	Keywords        string    `xml:"news:keywords,omitempty" json:"keywords,omitempty"`
}

// Alternate represents alternate versions of a page (mobile, print, etc.).
type Alternate struct {
	Media string `json:"media"`
	URL   string `json:"url"`
}

// Translation represents a translated version of a page.
type Translation struct {
	Language string `json:"language"`
	URL      string `json:"url"`
}

// Option is a functional option for configuring sitemap items.
type Option func(*Item)

// New creates a new sitemap with default options.
func New() *Sitemap {
	return &Sitemap{
		items: make([]Item, 0),
		opts: Options{
			MaxURLs: 50000,
		},
	}
}

// NewWithOptions creates a new sitemap with custom options.
func NewWithOptions(opts *Options) *Sitemap {
	if opts.MaxURLs <= 0 {
		opts.MaxURLs = 50000
	}

	items := make([]Item, 0)
	if opts.PreAllocate {
		items = make([]Item, 0, opts.MaxURLs)
	}

	return &Sitemap{
		items: items,
		opts:  *opts,
	}
}

// Add adds a URL to the sitemap with the specified parameters.
func (s *Sitemap) Add(loc string, lastMod time.Time, priority float64, changeFreq ChangeFreq, opts ...Option) error {
	if len(s.items) >= s.opts.MaxURLs {
		return fmt.Errorf("sitemap has reached maximum URL limit of %d", s.opts.MaxURLs)
	}

	if err := validateURL(loc); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if priority < 0.0 || priority > 1.0 {
		return fmt.Errorf("priority must be between 0.0 and 1.0, got %f", priority)
	}

	item := Item{
		URL:        loc,
		LastMod:    lastMod,
		Priority:   priority,
		ChangeFreq: changeFreq,
	}

	// Apply options
	for _, opt := range opts {
		opt(&item)
	}

	s.items = append(s.items, item)
	return nil
}

// AddItem adds a pre-configured item to the sitemap.
func (s *Sitemap) AddItem(item Item) error {
	if len(s.items) >= s.opts.MaxURLs {
		return fmt.Errorf("sitemap has reached maximum URL limit of %d", s.opts.MaxURLs)
	}

	if err := validateURL(item.URL); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if item.Priority < 0.0 || item.Priority > 1.0 {
		return fmt.Errorf("priority must be between 0.0 and 1.0, got %f", item.Priority)
	}

	s.items = append(s.items, item)
	return nil
}

// AddItems adds multiple items to the sitemap.
func (s *Sitemap) AddItems(items []Item) error {
	for _, item := range items {
		if err := s.AddItem(item); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the number of URLs in the sitemap.
func (s *Sitemap) Count() int {
	return len(s.items)
}

// Items returns all items in the sitemap.
func (s *Sitemap) Items() []Item {
	return s.items
}

// Clear removes all items from the sitemap.
func (s *Sitemap) Clear() {
	s.items = s.items[:0]
}

// WithTitle sets the title for a sitemap item.
func WithTitle(title string) Option {
	return func(item *Item) {
		item.Title = title
	}
}

// WithImages adds images to a sitemap item.
func WithImages(images []Image) Option {
	return func(item *Item) {
		item.Images = images
	}
}

// WithVideos adds videos to a sitemap item.
func WithVideos(videos []Video) Option {
	return func(item *Item) {
		item.Videos = videos
	}
}

// WithGoogleNews adds Google News metadata to a sitemap item.
func WithGoogleNews(news GoogleNews) Option {
	return func(item *Item) {
		item.News = &news
	}
}

// WithAlternates adds alternate versions to a sitemap item.
func WithAlternates(alternates []Alternate) Option {
	return func(item *Item) {
		item.Alternates = alternates
	}
}

// WithTranslations adds translations to a sitemap item.
func WithTranslations(translations []Translation) Option {
	return func(item *Item) {
		item.Langs = translations
	}
}

// validateURL validates that the URL is well-formed and absolute.
func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if !u.IsAbs() {
		return fmt.Errorf("URL must be absolute")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	return nil
}
