package apiv1

import (
	"ciao-admin/internal/server/router/endpoint"
	"log/slog"
	"net/http"
	"path"
)

const ApiV1Root = "/api/v1"

func GetApiV1Endpoints() []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:    path.Join(ApiV1Root, "token"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				slog.Info("get token triggered")
				w.WriteHeader(http.StatusOK)
			},
		},
	}
}
