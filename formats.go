package sitemap

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"time"
)

// TXT generates a plain text representation of the sitemap.
// Returns one URL per line.
func (s *Sitemap) TXT() ([]byte, error) {
	var buf bytes.Buffer
	for _, item := range s.items {
		buf.WriteString(item.URL)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

// HTML generates an HTML representation of the sitemap.
func (s *Sitemap) HTML() ([]byte, error) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sitemap</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .url-item { margin: 15px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .url { font-weight: bold; color: #0066cc; text-decoration: none; }
        .url:hover { text-decoration: underline; }
        .meta { color: #666; font-size: 0.9em; margin-top: 5px; }
        .images, .videos { margin-top: 10px; }
        .image, .video { margin: 5px 0; padding: 5px; background: #f9f9f9; border-radius: 3px; }
        .news { background: #fff3cd; padding: 10px; margin-top: 10px; border-radius: 5px; }
        h1 { color: #333; }
        .stats { background: #e9ecef; padding: 10px; border-radius: 5px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <h1>Sitemap</h1>
    <div class="stats">
        <strong>Total URLs:</strong> {{.Count}}
        {{if .LastGenerated}}<br><strong>Generated:</strong> {{.LastGenerated.Format "2006-01-02 15:04:05"}}{{end}}
    </div>
    {{range .Items}}
    <div class="url-item">
        <a href="{{.URL}}" class="url" target="_blank">{{.URL}}</a>
        <div class="meta">
            {{if .Priority}}<strong>Priority:</strong> {{printf "%.1f" .Priority}} | {{end}}
            {{if .ChangeFreq}}<strong>Change Frequency:</strong> {{.ChangeFreq}} | {{end}}
            {{if not .LastMod.IsZero}}<strong>Last Modified:</strong> {{.LastMod.Format "2006-01-02 15:04:05"}}{{end}}
        </div>
        {{if .Title}}<div class="meta"><strong>Title:</strong> {{.Title}}</div>{{end}}
        {{if .Images}}
        <div class="images">
            <strong>Images:</strong>
            {{range .Images}}
            <div class="image">
                {{if .Title}}<strong>{{.Title}}</strong><br>{{end}}
                <a href="{{.URL}}" target="_blank">{{.URL}}</a>
                {{if .Caption}}<br><em>{{.Caption}}</em>{{end}}
            </div>
            {{end}}
        </div>
        {{end}}
        {{if .Videos}}
        <div class="videos">
            <strong>Videos:</strong>
            {{range .Videos}}
            <div class="video">
                <strong>{{.Title}}</strong><br>
                {{.Description}}<br>
                {{if .ContentURL}}<a href="{{.ContentURL}}" target="_blank">Content URL</a> | {{end}}
                {{if .ThumbnailURL}}<a href="{{.ThumbnailURL}}" target="_blank">Thumbnail</a>{{end}}
                {{if .Duration}}<br><em>Duration: {{.Duration}} seconds</em>{{end}}
            </div>
            {{end}}
        </div>
        {{end}}
        {{if .News}}
        <div class="news">
            <strong>Google News:</strong><br>
            <strong>{{.News.Title}}</strong><br>
            Site: {{.News.SiteName}} | Language: {{.News.Language}}<br>
            Published: {{.News.PublicationDate.Format "2006-01-02 15:04:05"}}
            {{if .News.Keywords}}<br>Keywords: {{.News.Keywords}}{{end}}
        </div>
        {{end}}
    </div>
    {{end}}
</body>
</html>`

	t, err := template.New("sitemap").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	data := struct {
		Items         []Item
		Count         int
		LastGenerated *time.Time
	}{
		Items: s.items,
		Count: len(s.items),
		LastGenerated: func() *time.Time {
			now := time.Now()
			return &now
		}(),
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.Bytes(), err
}

// GoogleNews generates a Google News specific sitemap.
func (s *Sitemap) GoogleNews() ([]byte, error) {
	// Filter items that have Google News metadata
	var newsItems []Item
	for _, item := range s.items {
		if item.News != nil {
			newsItems = append(newsItems, item)
		}
	}

	// Create a temporary sitemap with only news items
	newsSitemap := &Sitemap{
		items: newsItems,
		opts:  s.opts,
	}

	return newsSitemap.XML()
}

// Mobile generates a mobile-specific sitemap.
func (s *Sitemap) Mobile() ([]byte, error) {
	urlset := struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		Mobile  string   `xml:"xmlns:mobile,attr"`
		URLs    []struct {
			URL    string `xml:"loc"`
			Mobile string `xml:"mobile:mobile,omitempty"`
		} `xml:"url"`
	}{
		Xmlns:  "http://www.sitemaps.org/schemas/sitemap/0.9",
		Mobile: "http://www.google.com/schemas/sitemap-mobile/1.0",
	}

	for _, item := range s.items {
		mobileURL := struct {
			URL    string `xml:"loc"`
			Mobile string `xml:"mobile:mobile,omitempty"`
		}{
			URL:    item.URL,
			Mobile: "", // This indicates it's a mobile page
		}
		urlset.URLs = append(urlset.URLs, mobileURL)
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

// JSON generates a JSON representation of the sitemap.
func (s *Sitemap) JSON() ([]byte, error) {
	return json.MarshalIndent(map[string]interface{}{
		"urls":  s.items,
		"count": len(s.items),
	}, "", "  ")
}
