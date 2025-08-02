module chi-example

go 1.23.0

toolchain go1.23.6

require (
	github.com/go-chi/chi/v5 v5.0.11
	github.com/rumendamyanov/go-sitemap v0.0.0-00010101000000-000000000000
)

replace github.com/rumendamyanov/go-sitemap => ../..
