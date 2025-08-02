# Chi Router Example

This example demonstrates how to use go-sitemap with the Chi router.

## Features Demonstrated

- Basic sitemap generation using Chi adapters
- Multiple output formats (XML, TXT, HTML)
- Custom sitemap handlers with standard HTTP handlers
- Dynamic content based on URL parameters
- Chi URL parameter extraction
- Image metadata support
- Comprehensive URL generation

## Running the Example

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

3. Access the sitemaps:
- **XML Sitemap**: http://localhost:8080/sitemap.xml
- **Text Sitemap**: http://localhost:8080/sitemap.txt
- **HTML Sitemap**: http://localhost:8080/sitemap.html
- **Custom Sitemap**: http://localhost:8080/custom-sitemap.xml
- **Query-based**: http://localhost:8080/custom-sitemap.xml?category=electronics
- **Dynamic Category**: http://localhost:8080/sitemap/books.xml

## Code Structure

- **Adapter Usage**: Uses `chiadapter.Sitemap()` for clean integration
- **Manual Handlers**: Shows how to build custom handlers for advanced use cases
- **URL Parameters**: Demonstrates Chi URL parameter extraction (`chi.URLParam`)
- **Standard HTTP**: Uses standard `http.Handler` interface for flexibility
- **Error Handling**: Proper error handling for production use

## Chi-Specific Features

- **URL Parameters**: Extract parameters from routes like `/sitemap/{category}.xml`
- **Middleware**: Compatible with Chi middleware stack
- **Standard HTTP**: Works with standard library patterns

## Learn More

- [Chi Router Documentation](https://github.com/go-chi/chi)
- [go-sitemap Wiki](../../wiki/)
- [Framework Integration Guide](../../wiki/Framework-Integration.md)
