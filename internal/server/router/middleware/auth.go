package middleware

import (
	"log/slog"
	"net/http"
)

type TokenPlace int
type ExtractTokenFunc func(*http.Request) (string, error)

const (
	Header = iota
	Cookie = 1
)

type AuthMiddleware struct {
	TokenPlace TokenPlace
}

func (auth *AuthMiddleware) Add(next http.HandlerFunc, onAuthFailed http.HandlerFunc, tokenExtractor ExtractTokenFunc) http.HandlerFunc {
	slog.Info("adding auth will check from header")
	switch auth.TokenPlace {
	case Header:
		return headerToken(next, onAuthFailed, tokenExtractor)
	case Cookie:
		return cookieToken(next, onAuthFailed, tokenExtractor)
	default:
		panic("unknow token place value")
	}
}

func headerToken(next http.HandlerFunc, onAuthFailed http.HandlerFunc, tokenExtractor ExtractTokenFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("jwt header is triggered")
		_, err := tokenExtractor(r)

		if err != nil {
			slog.Error("failed to extract token")
			onAuthFailed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func cookieToken(next http.HandlerFunc, onAuthFailed http.HandlerFunc, tokenExtractor ExtractTokenFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := tokenExtractor(r)

		if err != nil {
			slog.Error(err.Error())
			onAuthFailed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
