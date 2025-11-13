package endpoint

import "net/http"

type Endpoint struct {
	Path        string
	Methods     []string
	Secure      bool
	HandlerFunc http.HandlerFunc
}
