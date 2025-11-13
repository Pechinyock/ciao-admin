package middleware

import (
	"net/http"
)

type Middleware interface {
	Add(next http.Handler) http.Handler
}

type RequestMiddleware interface {
	Add(next http.HandlerFunc) http.HandlerFunc
}
