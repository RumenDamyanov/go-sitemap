# Echo Framework Example

This example demonstrates how to use go-sitemap with the Echo web framework.

## Features Demonstrated

- Basic sitemap generation using Echo adapters
- Multiple output formats (XML, TXT, HTML)
- Custom sitemap handlers with Echo context access
- Dynamic content based on request parameters
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

## Code Structure

- **Adapter Usage**: Uses `echoadapter.Sitemap()` for clean integration
- **Manual Handlers**: Shows how to build custom handlers for advanced use cases
- **Echo Context**: Demonstrates accessing request data for dynamic sitemaps
- **Error Handling**: Proper error handling for production use

## Learn More

- [Echo Framework Documentation](https://echo.labstack.com/)
- [go-sitemap Wiki](../../wiki/)
- [Framework Integration Guide](../../wiki/Framework-Integration.md)
