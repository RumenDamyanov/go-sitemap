package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"time"
)

// URLSet represents the root element of a sitemap XML.
type URLSet struct {
	XMLName xml.Name  `xml:"urlset"`
	Xmlns   string    `xml:"xmlns,attr"`
	Image   string    `xml:"xmlns:image,attr,omitempty"`
	Video   string    `xml:"xmlns:video,attr,omitempty"`
	News    string    `xml:"xmlns:news,attr,omitempty"`
	URLs    []XMLItem `xml:"url"`
}

// XMLItem represents a URL item in XML format.
type XMLItem struct {
	URL        string         `xml:"loc"`
	LastMod    string         `xml:"lastmod,omitempty"`
	ChangeFreq string         `xml:"changefreq,omitempty"`
	Priority   string         `xml:"priority,omitempty"`
	Images     []XMLImage     `xml:"image:image,omitempty"`
	Videos     []XMLVideo     `xml:"video:video,omitempty"`
	News       *XMLGoogleNews `xml:"news:news,omitempty"`
	Alternates []XMLAlternate `xml:"xhtml:link,omitempty"`
}

// XMLImage represents an image in XML format.
type XMLImage struct {
	URL     string `xml:"image:loc"`
	Title   string `xml:"image:title,omitempty"`
	Caption string `xml:"image:caption,omitempty"`
}

// XMLVideo represents a video in XML format.
type XMLVideo struct {
	ThumbnailURL string `xml:"video:thumbnail_loc"`
	Title        string `xml:"video:title"`
	Description  string `xml:"video:description"`
	ContentURL   string `xml:"video:content_loc,omitempty"`
	PlayerURL    string `xml:"video:player_loc,omitempty"`
	Duration     string `xml:"video:duration,omitempty"`
}

// XMLGoogleNews represents Google News metadata in XML format.
type XMLGoogleNews struct {
	Publication     XMLNewsPublication `xml:"news:publication"`
	PublicationDate string             `xml:"news:publication_date"`
	Title           string             `xml:"news:title"`
	Keywords        string             `xml:"news:keywords,omitempty"`
}

// XMLNewsPublication represents the news publication info.
type XMLNewsPublication struct {
	Name     string `xml:"news:name"`
	Language string `xml:"news:language"`
}

// XMLAlternate represents an alternate link in XML format.
type XMLAlternate struct {
	Rel      string `xml:"rel,attr"`
	Hreflang string `xml:"hreflang,attr,omitempty"`
	Media    string `xml:"media,attr,omitempty"`
	Href     string `xml:"href,attr"`
}

// XML generates the XML representation of the sitemap.
func (s *Sitemap) XML() ([]byte, error) {
	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]XMLItem, 0, len(s.items)),
	}

	// Check if we need namespace declarations
	hasImages, hasVideos, hasNews := false, false, false
	for _, item := range s.items {
		if len(item.Images) > 0 {
			hasImages = true
		}
		if len(item.Videos) > 0 {
			hasVideos = true
		}
		if item.News != nil {
			hasNews = true
		}
	}

	if hasImages {
		urlset.Image = "http://www.google.com/schemas/sitemap-image/1.1"
	}
	if hasVideos {
		urlset.Video = "http://www.google.com/schemas/sitemap-video/1.1"
	}
	if hasNews {
		urlset.News = "http://www.google.com/schemas/sitemap-news/0.9"
	}

	// Convert items to XML format
	for _, item := range s.items {
		xmlItem := XMLItem{
			URL: html.EscapeString(item.URL),
		}

		if !item.LastMod.IsZero() {
			xmlItem.LastMod = item.LastMod.Format(time.RFC3339)
		}

		if item.ChangeFreq != "" {
			xmlItem.ChangeFreq = string(item.ChangeFreq)
		}

		if item.Priority > 0 {
			xmlItem.Priority = formatPriority(item.Priority)
		}

		// Add images
		if len(item.Images) > 0 {
			xmlItem.Images = make([]XMLImage, len(item.Images))
			for i, img := range item.Images {
				xmlItem.Images[i] = XMLImage{
					URL:     html.EscapeString(img.URL),
					Title:   html.EscapeString(img.Title),
					Caption: html.EscapeString(img.Caption),
				}
			}
		}

		// Add videos
		if len(item.Videos) > 0 {
			xmlItem.Videos = make([]XMLVideo, len(item.Videos))
			for i, video := range item.Videos {
				xmlItem.Videos[i] = XMLVideo{
					ThumbnailURL: html.EscapeString(video.ThumbnailURL),
					Title:        html.EscapeString(video.Title),
					Description:  html.EscapeString(video.Description),
					ContentURL:   html.EscapeString(video.ContentURL),
					PlayerURL:    html.EscapeString(video.PlayerURL),
				}
				if video.Duration > 0 {
					xmlItem.Videos[i].Duration = formatDuration(video.Duration)
				}
			}
		}

		// Add Google News
		if item.News != nil {
			xmlItem.News = &XMLGoogleNews{
				Publication: XMLNewsPublication{
					Name:     html.EscapeString(item.News.SiteName),
					Language: item.News.Language,
				},
				PublicationDate: item.News.PublicationDate.Format(time.RFC3339),
				Title:           html.EscapeString(item.News.Title),
				Keywords:        html.EscapeString(item.News.Keywords),
			}
		}

		// Add alternates as xhtml:link elements
		if len(item.Alternates) > 0 || len(item.Langs) > 0 {
			for _, alt := range item.Alternates {
				xmlItem.Alternates = append(xmlItem.Alternates, XMLAlternate{
					Rel:   "alternate",
					Media: alt.Media,
					Href:  html.EscapeString(alt.URL),
				})
			}

			for _, lang := range item.Langs {
				xmlItem.Alternates = append(xmlItem.Alternates, XMLAlternate{
					Rel:      "alternate",
					Hreflang: lang.Language,
					Href:     html.EscapeString(lang.URL),
				})
			}
		}

		urlset.URLs = append(urlset.URLs, xmlItem)
	}

	// Generate XML
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteByte('\n')

	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlset); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// formatPriority formats priority value for XML output.
func formatPriority(priority float64) string {
	if priority == 1.0 {
		return "1.0"
	}
	if priority == 0.0 {
		return "0.0"
	}
	return fmt.Sprintf("%.1f", priority)
}

// formatDuration formats duration in seconds.
func formatDuration(duration int) string {
	return fmt.Sprintf("%d", duration)
}
