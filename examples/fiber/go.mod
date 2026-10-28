module fiber-example

go 1.26.10

require (
	github.com/gofiber/fiber/v2 v2.52.16
	go.rumenx.com/sitemap v0.0.0-00010101000000-000000000000
)

replace go.rumenx.com/sitemap => ../..

require (
	github.com/andybalholm/brotli v1.2.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.7 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
)
