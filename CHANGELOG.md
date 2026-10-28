# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-10-28

### Security

- Require Go 1.26.10 or newer. Go 1.25 and older are end of life and do not include the October 2026 toolchain fixes.
- Require `github.com/andybalholm/brotli` v1.2.7. Earlier versions panic on crafted Brotli input used by Gin and Fiber.
- Require `golang.org/x/net` v0.60.0 and `github.com/quic-go/quic-go` v0.59.1. Older copies are pulled in by Gin and are reachable from HTTP handlers.
