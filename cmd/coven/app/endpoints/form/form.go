package form

import (
	"ciao-admin/internal/server/router/endpoint"
	"log/slog"
	"net/http"
	"path"
)

const FormPrefix = "/form"

func GetFormEndpoints() []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:    path.Join(FormPrefix, "login"),
			Methods: []string{"POST"},
			Secure:  false,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				slog.Info("some one trying to login")
				w.WriteHeader(http.StatusUnauthorized)
			},
		},
	}
}
