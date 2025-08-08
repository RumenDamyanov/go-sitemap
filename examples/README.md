# go-sitemap Examples

This directory contains working examples demonstrating how to use go-sitemap with different Go web frameworks and the standard library.

## Available Examples

### [Standard net/http](nethttp/)
Basic example using Go's standard `net/http` library.
- ✅ Minimal dependencies
- ✅ Manual XML handling
- ✅ Standard library only

### [Gin Framework](gin/)
Integration with the Gin web framework.
- ✅ Gin adapters (`ginadapter.Sitemap()`)
- ✅ Multiple output formats
- ✅ Custom handlers with Gin context
- ✅ Middleware integration

### [Echo Framework](echo/)
Integration with the Echo web framework.
- ✅ Echo adapters (`echoadapter.Sitemap()`)
- ✅ Multiple output formats
- ✅ Custom handlers with Echo context
- ✅ Middleware support

### [Fiber Framework](fiber/)
Integration with the Fiber web framework.
- ✅ Fiber adapters (`fiberadapter.Sitemap()`)
- ✅ Multiple output formats
- ✅ Custom handlers with Fiber context
- ✅ Fast performance

### [Chi Router](chi/)
Integration with the Chi router.
- ✅ Chi adapters (`chiadapter.Sitemap()`)
- ✅ URL parameter extraction
- ✅ Standard HTTP handler compatibility
- ✅ Lightweight routing

## Running Examples

Each example can be run independently:

```bash
# Choose any framework
cd gin/           # or echo/, fiber/, chi/, nethttp/
go mod tidy       # Install dependencies
go run main.go    # Start the server
```

Then visit:

- **XML Sitemap**: <http://localhost:8080/sitemap.xml>
- **Text Sitemap**: <http://localhost:8080/sitemap.txt>
- **HTML Sitemap**: <http://localhost:8080/sitemap.html>

## Common Features

All examples demonstrate:

- **Basic sitemap generation** with URLs, priorities, and change frequencies
- **Multiple output formats** (XML, TXT, HTML)
- **Dynamic content** based on request parameters
- **Image metadata** support
- **Error handling** for production use
- **Custom handlers** for advanced use cases

## Choosing a Framework

| Framework | Best For | Performance | Learning Curve |
|-----------|----------|-------------|----------------|
| **net/http** | Minimal deps, learning | Good | Easy |
| **Gin** | Rapid development | Very Good | Easy |
| **Echo** | Balanced features | Very Good | Easy |
| **Fiber** | High performance | Excellent | Medium |
| **Chi** | Lightweight routing | Good | Easy |

## Learn More

- 📚 [go-sitemap Wiki](../wiki/)
- 🚀 [Quick Start Guide](../wiki/Quick-Start.md)
- 🔧 [Framework Integration](../wiki/Framework-Integration.md)
- 🎯 [Best Practices](../wiki/Best-Practices.md)

## Contributing

Found an issue or want to improve an example? Please see our [Contributing Guidelines](../CONTRIBUTING.md).
