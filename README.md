# go-sitemap

[![CI](https://github.com/rumendamyanov/go-sitemap/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-sitemap/actions/workflows/ci.yml)
![CodeQL](https://github.com/rumendamyanov/go-sitemap/actions/workflows/github-code-scanning/codeql/badge.svg)
![Dependabot](https://github.com/rumendamyanov/go-sitemap/actions/workflows/dependabot/dependabot-updates/badge.svg)
[![codecov](https://codecov.io/gh/rumendamyanov/go-sitemap/branch/master/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-sitemap)
[![Go Report Card](https://goreportcard.com/badge/go.rumenx.com/sitemap?)](https://goreportcard.com/report/go.rumenx.com/sitemap)
[![Go Reference](https://pkg.go.dev/badge/github.com/rumendamyanov/go-sitemap.svg)](https://pkg.go.dev/github.com/rumendamyanov/go-sitemap)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rumendamyanov/go-sitemap/blob/master/LICENSE.md)

A framework-agnostic Go module for generating sitemaps in XML, TXT, HTML, and Google News formats. Inspired by [php-sitemap](https://github.com/RumenDamyanov/php-sitemap), this package works seamlessly with any Go web framework including Gin, Echo, Fiber, Chi, and standard net/http.

## Features

- **Framework-agnostic**: Use with Gin, Echo, Fiber, Chi, or standard net/http
- **Multiple formats**: XML, TXT, HTML, Google News, mobile sitemaps
- **Rich content**: Supports images, videos, translations, alternates, Google News
- **Modern Go**: Type-safe, extensible, and robust (Go 1.22+)
- **High test coverage**: Comprehensive test suite with CI/CD integration
- **Easy integration**: Simple API, drop-in for handlers/middleware
- **Extensible**: Adapters for popular Go web frameworks
- **Production ready**: Used in production environments

## Quick Links

- 📖 [Installation](#installation)
- 🚀 [Usage Examples](#usage)
- 🔧 [Framework Adapters](#framework-adapters)
- 📚 [Documentation Wiki](wiki/)
- 🧪 [Testing & Development](#testing--development)
- 🤝 [Contributing](CONTRIBUTING.md)
- 🔒 [Security Policy](SECURITY.md)
- 💝 [Support & Funding](FUNDING.md)
- 📄 [License](#license)

## Installation

```bash
go get go.rumenx.com/sitemap
```

## Usage

### Basic Example (net/http)

```go
package main

import (
    "net/http"
    "time"

    "go.rumenx.com/sitemap"
)

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
    sm := sitemap.New()

    // Add URLs
    sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
    sm.Add("https://example.com/about", time.Now(), 0.8, sitemap.Monthly,
        sitemap.WithImages([]sitemap.Image{
            {URL: "https://example.com/img/about.jpg", Title: "About Us"},
        }),
    )

    // Render as XML
    w.Header().Set("Content-Type", "application/xml")
    xml, err := sm.XML()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(xml)
}

func main() {
    http.HandleFunc("/sitemap.xml", sitemapHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Advanced Features

```go
sm := sitemap.New()

// Add with all supported fields
sm.Add(
    "https://example.com/news",
    time.Now(),
    0.8,
    sitemap.Daily,
    sitemap.WithImages([]sitemap.Image{
        {URL: "https://example.com/img/news.jpg", Title: "News Image"},
    }),
    sitemap.WithTitle("News Article"),
    sitemap.WithTranslations([]sitemap.Translation{
        {Language: "fr", URL: "https://example.com/fr/news"},
    }),
    sitemap.WithVideos([]sitemap.Video{
        {Title: "News Video", Description: "Video description"},
    }),
    sitemap.WithGoogleNews(sitemap.GoogleNews{
        SiteName: "Example News",
        Language: "en",
        PublicationDate: time.Now(),
    }),
    sitemap.WithAlternates([]sitemap.Alternate{
        {Media: "print", URL: "https://example.com/news-print"},
    }),
)

// Multiple output formats
xmlData, _ := sm.XML()        // Standard XML sitemap
txtData, _ := sm.TXT()        // Plain text format
htmlData, _ := sm.HTML()      // HTML format
newsData, _ := sm.GoogleNews() // Google News sitemap
```

## Framework Adapters

### Gin Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.rumenx.com/sitemap/adapters/gin"
)

func main() {
    r := gin.Default()

    r.GET("/sitemap.xml", ginadapter.Sitemap(func() *sitemap.Sitemap {
        sm := sitemap.New()
        sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
        return sm
    }))

    r.Run(":8080")
}
```

### Echo Example

```go
package main

import (
    "github.com/labstack/echo/v4"
    "go.rumenx.com/sitemap/adapters/echo"
)

func main() {
    e := echo.New()

    e.GET("/sitemap.xml", echoadapter.Sitemap(func() *sitemap.Sitemap {
        sm := sitemap.New()
        sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
        return sm
    }))

    e.Start(":8080")
}
```

### Fiber Example

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "go.rumenx.com/sitemap/adapters/fiber"
)

func main() {
    app := fiber.New()

    app.Get("/sitemap.xml", fiberadapter.Sitemap(func() *sitemap.Sitemap {
        sm := sitemap.New()
        sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
        return sm
    }))

    app.Listen(":8080")
}
```

### Chi Example

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "go.rumenx.com/sitemap/adapters/chi"
)

func main() {
    r := chi.NewRouter()

    r.Get("/sitemap.xml", chiadapter.Sitemap(func() *sitemap.Sitemap {
        sm := sitemap.New()
        sm.Add("https://example.com/", time.Now(), 1.0, sitemap.Daily)
        return sm
    }))

    http.ListenAndServe(":8080", r)
}
```

## Multiple Methods for Adding Items

### add() vs addItem()

You can add sitemap entries using either the `Add()` or `AddItem()` methods:

**Add() — Simple, type-safe, one-at-a-time:**

```go
// Recommended for most use cases
sm.Add(
    "https://example.com/",
    time.Now(),
    1.0,
    sitemap.Daily,
    sitemap.WithImages([]sitemap.Image{
        {URL: "https://example.com/img.jpg", Title: "Image"},
    }),
    sitemap.WithTitle("Homepage"),
)
```

**AddItem() — Advanced, struct-based, supports batch:**

```go
// Add a single item with a struct
sm.AddItem(sitemap.Item{
    URL:      "https://example.com/about",
    LastMod:  time.Now(),
    Priority: 0.8,
    ChangeFreq: sitemap.Monthly,
    Title:    "About Us",
    Images:   []sitemap.Image{{URL: "https://example.com/img/about.jpg", Title: "About"}},
})

// Add multiple items at once (batch add)
sm.AddItems([]sitemap.Item{
    {URL: "https://example.com/page1", Title: "Page 1"},
    {URL: "https://example.com/page2", Title: "Page 2"},
})
```

## Documentation

For comprehensive documentation and examples:

- 📚 [Quick Start Guide](wiki/Quick-Start.md) - Get up and running quickly
- 🔧 [Basic Usage](wiki/Basic-Usage.md) - Core functionality and examples
- 🚀 [Advanced Usage](wiki/Advanced-Usage.md) - Advanced features and customization
- 🔌 [Framework Integration](wiki/Framework-Integration.md) - Integration with popular frameworks
- 🎯 [Best Practices](wiki/Best-Practices.md) - Performance tips and recommendations
- 🤝 [Contributing Guidelines](CONTRIBUTING.md) - How to contribute to this project
- 🔒 [Security Policy](SECURITY.md) - Security guidelines and vulnerability reporting
- 💝 [Funding & Support](FUNDING.md) - Support and sponsorship information

## Testing & Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Code Quality

```bash
# Run static analysis
go vet ./...

# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Development setup
- Coding standards
- Testing requirements
- Pull request process

## Security

If you discover a security vulnerability, please review our [Security Policy](SECURITY.md) for responsible disclosure guidelines.

## Support

If you find this package helpful, consider:

- ⭐ Starring the repository
- 💝 [Supporting development](FUNDING.md)
- 🐛 [Reporting issues](https://github.com/rumendamyanov/go-sitemap/issues)
- 🤝 [Contributing improvements](CONTRIBUTING.md)

## License

[MIT License](LICENSE.md)
