package middleware

import (
	"log/slog"
	"net/http"
)

type TokenPlace int

const (
	Header = iota
	Cookie = 1
)

type AuthMiddleware struct {
	TokenPlace TokenPlace
}

func (auth *AuthMiddleware) Add(next http.HandlerFunc) http.HandlerFunc {
	slog.Info("adding auth will check from header")
	switch auth.TokenPlace {
	case Header:
		return headerToken(next)
	case Cookie:
		return cookieToken(next)
	default:
		panic("unknow token place value")
	}
}

func headerToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("jwt header is triggered")
		next.ServeHTTP(w, r)
	})
}

func cookieToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("jwt cookie is triggered")
		next.ServeHTTP(w, r)
	})
}
