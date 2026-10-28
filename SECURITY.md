# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | ✅                |
| < 1.0   | ❌                |

Build with **Go 1.26.10+** or **1.27.3+**. Go 1.25 and older are end of life and do not receive the October 2026 security releases. Those fixes are in the toolchain (`crypto/x509`, `net/http`, `html/template`), not in this module. The `go` line asks the `go` command for at least Go 1.26.10 when `GOTOOLCHAIN=auto`, which is the default. `GOTOOLCHAIN=local` fails the build until that toolchain is installed.

- [CVE-2026-111943](https://pkg.go.dev/vuln/GO-2026-6809) (`crypto/x509`): untrustworthy eMudhra intermediates are rejected. Fixed in Go 1.26.10 and 1.27.3 ([announcement](https://groups.google.com/g/golang-announce/c/66pojueaQek/m/tmWX2Pq9AgAJ)).
- Go 1.26.9 and 1.27.2 also fixed HTTP/2 server crashes, trailer memory use, and `html/template` context tracking ([announcement](https://groups.google.com/g/golang-announce/c/U2fTuyDJznI)).
- Gin and Fiber decompress Brotli with `github.com/andybalholm/brotli`. Versions before v1.2.7 panic on crafted input ([CVE-2026-111497](https://www.wiz.io/blog/cve-2026-111497-brotli-rce-node-js)). This module requires v1.2.7.

## Reporting a Vulnerability

We take the security of go-sitemap seriously. If you discover a security vulnerability, please follow these steps:

### 1. Do NOT open a public GitHub issue

Security vulnerabilities should be reported privately to allow us to fix them before they are publicly disclosed.

### 2. Send a report via email

Email your findings to: **contact@rumenx.com**

Please include:
- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)

### 3. Response timeline

- **Initial response**: Within 48 hours
- **Status update**: Within 7 days
- **Fix timeline**: Depends on severity, typically within 30 days

### 4. Disclosure policy

- We will acknowledge receipt of your vulnerability report
- We will investigate and validate the issue
- We will work on a fix and release it as soon as possible
- We will coordinate with you on the disclosure timeline
- We will publicly credit you for the discovery (unless you prefer to remain anonymous)

## Security best practices

When using go-sitemap:

1. **Validate input URLs** - Always validate and sanitize URLs before adding them to sitemaps
2. **Rate limiting** - Implement rate limiting for sitemap generation endpoints
3. **Memory management** - Be mindful of memory usage when generating large sitemaps
4. **Access control** - Secure your sitemap generation endpoints appropriately
5. **Regular updates** - Keep the package updated to the latest version

## Security considerations

- This package generates XML output, ensure proper escaping of user input
- Large sitemaps can consume significant memory, implement appropriate limits
- Be cautious when accepting URLs from external sources

## Contact

For any security-related questions or concerns, please contact:

**Email:** contact@rumenx.com
**Maintainer:** Rumen Damyanov

Thank you for helping keep go-sitemap and our users safe!
