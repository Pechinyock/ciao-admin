package app

import (
	"ciao-admin/cmd/coven/app/configuration"
	ui "ciao-admin/cmd/coven/app/endpoints/UI"
	"ciao-admin/cmd/coven/app/endpoints/form"
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/server/router/middleware"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

func registerFromEndpoints(router *http.ServeMux) error {
	if router == nil {
		return errors.New("trying to add endpoint on to nil router")
	}

	formRequests := form.GetFormEndpoints()
	slog.Info("register forms endpoints", "total forms endpoints", len(formRequests))

	fromsAuth := middleware.AuthMiddleware{
		TokenPlace: middleware.Cookie,
	}

	for _, route := range formRequests {
		if route.Secure {
			router.HandleFunc(route.Path, fromsAuth.Add(route.HandlerFunc,
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				},
				extractTokenFromCookie))
		} else {
			router.HandleFunc(route.Path, route.HandlerFunc)
		}
	}
	return nil
}

func registerUIEndpoints(router *http.ServeMux, uiBundle *webui.WebUIBundle) error {
	if router == nil {
		return errors.New("trying to add endpoint on to nil router")
	}
	if uiBundle == nil {
		return errors.New("trying to add ui endpoints with nil ui bundle")
	}

	uiEndpoints := ui.GetUIEndpoints(uiBundle)
	slog.Info("register forms endpoints", "total forms endpoints", len(uiEndpoints))

	uiAuth := middleware.AuthMiddleware{
		TokenPlace: middleware.Cookie,
	}

	for _, route := range uiEndpoints {
		if route.Secure {
			router.HandleFunc(route.Path, uiAuth.Add(route.HandlerFunc,
				redirectToLoginScreen,
				extractTokenFromCookie,
			))
		} else {
			router.HandleFunc(route.Path, route.HandlerFunc)
		}
	}
	return nil
}

func registerWebUIStaticFilesEndpoints(router *http.ServeMux, config *webui.WebUIBundleConfig) error {
	if router == nil {
		return errors.New("failed to add static files route provided router is nil")
	}
	err := validateConfig(config)
	if err != nil {
		return err
	}
	staticFilesDirPath := filepath.Join(config.RootPath, config.StaticFilesDirName)
	fs := http.FileServer(http.Dir(staticFilesDirPath))
	routeName := fmt.Sprintf("/%s/", config.StaticFilesRootRouteName)
	router.Handle(routeName, http.StripPrefix(routeName, fs))
	return nil
}

func registerFileShareEndpoints(router *http.ServeMux, config *configuration.FileServerConfig) error {
	if router == nil {
		return errors.New("failed to add file share routes, provided router is nil")
	}
	if config == nil || len(config.ShareDirConfigs) == 0 {
		return errors.New("failed to add file share routes, provided config is nil")
	}
	handlerSetter := func(routeName, path, source string) http.Handler {
		fs := http.FileServer(http.Dir(path))

		switch strings.ToLower(source) {
		case "header":
			panic("not implemented")
		case "cookie":
			panic("not implemented")
		case "none", "":
			{
				formated := fmt.Sprintf("/%s/", routeName)
				router.Handle(formated, http.StripPrefix(formated, fs))
				slog.Info("succesfly added directory as files server",
					"directory physical path", path,
					"route name", routeName,
				)
				return fs
			}
		default:
			panic(fmt.Sprintf("unknown token place for file server \n directory name: %s\n route name: %s", path, routeName))
		}
	}
	for _, e := range config.ShareDirConfigs {
		handlerSetter(e.RouteName, e.DirPath, e.TokenSource)
	}
	return nil
}

func redirectToLoginScreen(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	return "", nil
}
