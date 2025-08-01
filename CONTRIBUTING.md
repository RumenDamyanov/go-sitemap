# Contributing to go-sitemap

Thank you for your interest in contributing to go-sitemap! We welcome contributions from the community and are grateful for your help in making this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- Basic understanding of Go modules and packages

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/go-sitemap.git
cd go-sitemap
```

3. Add the original repository as upstream:

```bash
git remote add upstream https://github.com/rumendamyanov/go-sitemap.git
```

## Development Setup

1. Install dependencies:

```bash
go mod download
```

2. Run tests to ensure everything works:

```bash
go test ./...
```

3. Run tests with coverage:

```bash
go test -cover ./...
```

## Contributing Guidelines

### Types of Contributions

We welcome several types of contributions:

- 🐛 **Bug fixes** - Fix issues and problems
- ✨ **New features** - Add new functionality
- 📝 **Documentation** - Improve or add documentation
- 🧪 **Tests** - Add or improve test coverage
- 🔧 **Framework adapters** - Add support for new frameworks
- 🎨 **Code improvements** - Refactor and optimize existing code

### Before You Start

1. **Check existing issues** - Look for existing issues or discussions
2. **Create an issue** - For significant changes, create an issue first
3. **Discuss the approach** - Get feedback on your proposed solution
4. **Keep it focused** - One feature/fix per pull request

## Pull Request Process

### 1. Create a Branch

Create a descriptive branch name:

```bash
git checkout -b feature/add-chi-adapter
git checkout -b fix/memory-leak-large-sitemaps
git checkout -b docs/improve-readme-examples
```

### 2. Make Your Changes

- Write clean, readable code
- Follow the coding standards (see below)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Check test coverage
go test -cover ./...

# Run specific tests
go test ./adapters/gin/

# Check formatting
go fmt ./...

# Run static analysis
go vet ./...
```

### 4. Commit Your Changes

Use clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add Chi framework adapter

- Add ChiAdapter for Chi router integration
- Include comprehensive tests and examples
- Update documentation with Chi usage examples"
```

### 5. Push and Create Pull Request

```bash
git push origin your-branch-name
```

Then create a pull request on GitHub with:

- Clear title and description
- Reference any related issues
- Include screenshots/examples if applicable

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Follow Go naming conventions
- Keep functions focused and concise

### Example Code Style

```go
// Package-level comment
package sitemap

// Image represents an image reference in a sitemap entry.
// It follows the Google Images sitemap specification.
type Image struct {
    // URL is the absolute URL of the image.
    URL string `xml:"image:loc" json:"url"`

    // Title provides a short description of the image.
    Title string `xml:"image:title,omitempty" json:"title,omitempty"`

    // Caption describes the image content.
    Caption string `xml:"image:caption,omitempty" json:"caption,omitempty"`
}

// Add adds a new URL to the sitemap with the specified parameters.
// It returns an error if the URL is invalid or if the sitemap is sealed.
func (s *Sitemap) Add(url string, lastMod time.Time, priority float64, changeFreq ChangeFreq, opts ...Option) error {
    if err := s.validateURL(url); err != nil {
        return fmt.Errorf("invalid URL: %w", err)
    }

    // Implementation...
    return nil
}
```

### Documentation Standards

- All exported functions must have comments
- Include examples in documentation
- Update README for significant changes
- Add wiki documentation for complex features

## Testing Requirements

### Test Coverage

- Aim for high test coverage (>90%)
- Test both success and error cases
- Include edge cases and boundary conditions
- Test framework adapters thoroughly

### Test Structure

```go
func TestSitemap_Add(t *testing.T) {
    tests := []struct {
        name        string
        url         string
        priority    float64
        expectError bool
    }{
        {
            name:        "valid URL",
            url:         "https://example.com/",
            priority:    1.0,
            expectError: false,
        },
        {
            name:        "invalid URL",
            url:         "not-a-url",
            priority:    1.0,
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := New()
            err := s.Add(tt.url, time.Now(), tt.priority, Daily)

            if tt.expectError && err == nil {
                t.Error("expected error, got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### Testing Framework Adapters

- Test adapter integration with framework
- Mock framework contexts appropriately
- Test error handling and edge cases
- Include integration examples

## Documentation

### README Updates

When adding new features:

- Update the feature list
- Add usage examples
- Update the table of contents
- Add links to wiki documentation

### Wiki Documentation

For comprehensive features, add wiki documentation:

- Create detailed guides
- Include multiple examples
- Explain best practices
- Add troubleshooting sections

### Code Comments

```go
// Package sitemap provides functionality for generating XML sitemaps
// following the sitemaps.org protocol. It supports standard sitemaps,
// Google News sitemaps, image sitemaps, and video sitemaps.
//
// The package is designed to be framework-agnostic and includes
// adapters for popular Go web frameworks.
package sitemap
```

## Community

### Getting Help

- Create an issue for bugs or questions
- Join discussions in existing issues
- Follow the project for updates

### Communication

- Be respectful and constructive
- Provide detailed information in issues
- Be patient with response times
- Help other community members

## Recognition

Contributors are recognized in:

- Project README
- Release notes
- Git commit history
- Special thanks in documentation

## Questions?

If you have questions about contributing, please:

1. Check existing documentation
2. Search existing issues
3. Create a new issue with the "question" label
4. Contact the maintainer: [contact@rumenx.com](mailto:contact@rumenx.com)

Thank you for contributing to go-sitemap! 🚀
