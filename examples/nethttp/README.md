# Standard net/http Example

This example demonstrates how to use go-sitemap with the standard Go net/http library.

## Features Demonstrated

- Basic sitemap generation using standard HTTP handlers
- Manual XML response handling
- Custom sitemap logic without framework adapters
- Direct integration with `http.ServeMux`
- Error handling and content type setting

## Running the Example

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

3. Access the sitemap:

- **XML Sitemap**: <http://localhost:8080/sitemap.xml>

## Code Structure

- **Standard Library**: Uses only Go's built-in `net/http` package
- **Manual Handling**: Shows how to manually create and serve sitemaps
- **Direct Integration**: No framework dependencies
- **Error Handling**: Proper HTTP error responses

## When to Use

This approach is ideal when:

- You want minimal dependencies
- Using the standard library only
- Building a custom web framework
- Learning the underlying mechanics

## Learn More

- [Go net/http Documentation](https://pkg.go.dev/net/http)
- [go-sitemap Wiki](../../wiki/)
- [Basic Usage Guide](../../wiki/Basic-Usage.md)
