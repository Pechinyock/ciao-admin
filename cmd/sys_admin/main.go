package main

import "fmt"

var (
	Version    = "dev"
	GitShorSha = "unknown"
)

func main() {
	fmt.Printf("sys admin version: %s", Version)
}
