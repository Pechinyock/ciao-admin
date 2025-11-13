package app

import (
	"ciao-admin/cmd/admin/app/configuration"
	apiv1 "ciao-admin/cmd/admin/app/endpoints/api_v1"
	"ciao-admin/cmd/admin/app/endpoints/form"
	"ciao-admin/internal/UI/webui"
	"ciao-admin/internal/server/router/middleware"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
)

func registerApiEndpoints(router *http.ServeMux) error {
	if router == nil {
		return errors.New("trying to add endpoint on to nil router")
	}

	apiv1 := apiv1.GetApiV1Endpoints()
	apiAuth := middleware.AuthMiddleware{
		TokenPlace: middleware.Header,
	}

	slog.Info("api endpoints are enabled", "api version", 1, "total api endpoints", len(apiv1))
	for _, route := range apiv1 {
		if route.Secure {
			router.HandleFunc(route.Path, apiAuth.Add(route.HandlerFunc))
		} else {
			router.HandleFunc(route.Path, route.HandlerFunc)
		}
	}

	return nil
}

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
			router.HandleFunc(route.Path, fromsAuth.Add(route.HandlerFunc))
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
