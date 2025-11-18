package ui

import (
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/server/router/endpoint"
	"net/http"
	"path"
)

const UIPrefix = "/ui"

func GetUIEndpoints(uiBundle *webui.WebUIBundle) []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:    "/",
			Methods: []string{"GET"},
			Secure:  true,
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
		{
			Path:    "/login",
			Methods: []string{"GET"},
			Secure:  false,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				uiBundle.Render("login_screen", w, nil)
			},
		},
		{
			Path:    path.Join(UIPrefix, "main-menu"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				uiBundle.Render("menu", w, nil)
			},
		},
		{
			Path:    path.Join(UIPrefix, "coven"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				err := uiBundle.Render("coven", w, nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		},
		{
			Path:    path.Join(UIPrefix, "create-card"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				err := uiBundle.Render("create_card", w, nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		},
	}
}
