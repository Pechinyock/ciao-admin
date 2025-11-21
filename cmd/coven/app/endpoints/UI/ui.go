package ui

import (
	"ciao-admin/cmd/coven/app/cards"
	"ciao-admin/cmd/coven/app/projection"
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/server/router/endpoint"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
)

var UIBundle *webui.WebUIBundle

const UIPrefix = "/ui"

func GetUIEndpoints(uiBundle *webui.WebUIBundle) []endpoint.Endpoint {
	return []endpoint.Endpoint{
		{
			Path:    "/",
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				uiBundle.Render("main", w, nil)
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
			Path:    path.Join(UIPrefix, "modal-window"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				err := uiBundle.Render("modal_window", w, nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			},
		},
		{
			Path:    path.Join(UIPrefix, "modal-body", "{modalName}"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				modalName := r.PathValue("modalName")
				if modalName == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				templName := strings.ReplaceAll(modalName, "-", "_")
				renderFunc, defined := modalsMap[templName]
				if !defined {
					slog.Error("modal winodw not found")
					w.WriteHeader(http.StatusNotFound)
					return
				}
				renderFunc(templName, w)
			},
		},
		{
			Path:    path.Join(UIPrefix, "create-card-from"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				cardType := r.FormValue("create-card-group")
				cardTypeLower := strings.ToLower(cardType)
				_, isExists := cards.CardTypes[cardTypeLower]
				if !isExists {
					slog.Error("failed to render card create form: unknown type", "provided type", cardType)
					w.WriteHeader(http.StatusNotFound)
				}
				slog.Info("create card form", "asdf", cardType)

				switch cardTypeLower {
				case "characters":
					uiBundle.Render("create_character", w, nil)
				case "spells":
					w.Write([]byte("unimplemented spells"))
				case "secrets":
					w.Write([]byte("unimplemented secrets"))
				case "curses":
					w.Write([]byte("unimplemented curses"))
				case "ingredients":
					w.Write([]byte("unimplemented ingredients"))
				case "potions":
					w.Write([]byte("unimplemented potions"))
				default:
					slog.Error("unknown card type", "type", cardTypeLower)
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
		},
		{
			Path:    path.Join(UIPrefix, "image-pool", "{poolName}"),
			Methods: []string{"GET"},
			Secure:  true,
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				poolName := r.PathValue("poolName")
				if poolName == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				poolNameLower := strings.ToLower(poolName)
				_, defined := cards.CardTypes[poolNameLower]
				if !defined {
					slog.Error("unknown type of image pool")
					w.WriteHeader(http.StatusNotFound)
					return
				}
				prewiewData, err := loadImagesPrewiewData(poolName)
				if err != nil {
					slog.Error("failed to load image pool view", "error message", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if prewiewData == nil {
					slog.Info("there's no images to display")
					w.WriteHeader(http.StatusNoContent)
					w.Write([]byte("Картинки данного типа еще не загружены"))
					return
				}
				poolView := projection.ImageViewProj{
					BasePath:  baseUriPath,
					FileGroup: poolNameLower,
					Images:    prewiewData,
				}

				err = uiBundle.Render("select_image", w, poolView)
				if err != nil {
					slog.Error("failed to render image pool", "group", poolName,
						"error message", err.Error(),
					)
					errMsg := fmt.Sprintf("Не удалось отобразить картинки для группы: %s", poolName)
					uiBundle.Render("alert", w, projection.AlertProj{
						Message: errMsg,
						Type:    "danger",
					})
					return
				}
			},
		},
	}
}
