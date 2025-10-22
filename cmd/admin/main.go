package main

import "fmt"

var (
	Version    = "dev"
	GitShorSha = "unknown"
)

func main() {
	fmt.Printf("admin version: %s", Version)
}
