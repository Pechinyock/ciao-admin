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
			Path:    path.Join(FormPrefix, "token"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				slog.Info("get token form triggered")
				w.WriteHeader(http.StatusOK)
			},
		},
	}
}
