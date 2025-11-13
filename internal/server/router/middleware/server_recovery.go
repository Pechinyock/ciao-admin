package middleware

import "net/http"

type ServerRecovery struct {
}

func (sr *ServerRecovery) Add(request http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		request.ServeHTTP(w, r)
	})
}
