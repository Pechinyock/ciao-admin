package ui

import (
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/server/router/endpoint"
	"net/http"
)

const UIPrefix = "/ui"

func GetUIEndpoints(uiBundle *webui.WebUIBundle) []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:    "/",
			Methods: []string{"GET"},
			Secure:  false,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				if uiBundle == nil {
					panic("ui bundle nil")
				}
				person := struct {
					Name string
				}{
					Name: "Mi-si pisi",
				}
				uiBundle.Render("main", w, person)
			},
		},
	}
}
