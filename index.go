package sitemap

import (
	"bytes"
	"encoding/xml"
	"html"
	"time"
)

// Index represents a sitemap index that references multiple sitemaps.
type Index struct {
	sitemaps []IndexItem
}

// IndexItem represents a single sitemap reference in the index.
type IndexItem struct {
	URL     string    `xml:"loc"`
	LastMod time.Time `xml:"lastmod,omitempty"`
}

// IndexURLSet represents the root element of a sitemap index XML.
type IndexURLSet struct {
	XMLName  xml.Name       `xml:"sitemapindex"`
	Xmlns    string         `xml:"xmlns,attr"`
	Sitemaps []IndexXMLItem `xml:"sitemap"`
}

// IndexXMLItem represents a sitemap item in XML format for the index.
type IndexXMLItem struct {
	URL     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// NewIndex creates a new sitemap index.
func NewIndex() *Index {
	return &Index{
		sitemaps: make([]IndexItem, 0),
	}
}

// Add adds a sitemap URL to the index.
func (idx *Index) Add(url string, lastMod time.Time) error {
	if err := validateURL(url); err != nil {
		return err
	}

	idx.sitemaps = append(idx.sitemaps, IndexItem{
		URL:     url,
		LastMod: lastMod,
	})

	return nil
}

// Count returns the number of sitemaps in the index.
func (idx *Index) Count() int {
	return len(idx.sitemaps)
}

// XML generates the XML representation of the sitemap index.
func (idx *Index) XML() ([]byte, error) {
	urlset := IndexURLSet{
		Xmlns:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: make([]IndexXMLItem, len(idx.sitemaps)),
	}

	for i, sitemap := range idx.sitemaps {
		urlset.Sitemaps[i] = IndexXMLItem{
			URL: html.EscapeString(sitemap.URL),
		}

		if !sitemap.LastMod.IsZero() {
			urlset.Sitemaps[i].LastMod = sitemap.LastMod.Format(time.RFC3339)
		}
	}

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
