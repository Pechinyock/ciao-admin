package webui

import (
	"ciao-admin/internal/utils"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	rootTemplateName = "root"
)

type WebUIBundle struct {
	config    *WebUIBundleConfig
	rootTempl *template.Template
}

func New(config *WebUIBundleConfig) (*WebUIBundle, error) {
	if config.StaticFilesDirName == "" {
		return nil, errors.New("static files dir name was not provided")
	}
	bundle := WebUIBundle{
		config: config,
	}

	allTemplates := config.TemplatePaths

	if len(allTemplates) == 0 {
		return nil, errors.New("failed to load web bundle theres no templates were provided")
	}

	bundle.loadTemplates(allTemplates)
	return &bundle, nil
}

func (b *WebUIBundle) Render(templName string, writter io.Writer, data any) error {
	return b.rootTempl.ExecuteTemplate(writter, templName, data)
}

func (b *WebUIBundle) loadTemplates(templatesPaths ...[]string) {
	if len(templatesPaths) == 0 {
		slog.Error("failed to load laytout tempates, list of templates is empty")
		return
	}

	allTemplateFilePaths := make([]string, 0)
	for _, tmplPaths := range templatesPaths {
		for _, templPath := range tmplPaths {
			fullPath := path.Join(b.config.RootPath, templPath)
			if utils.IsGlob(fullPath) {
				paths, err := filepath.Glob(fullPath)
				if len(paths) == 0 {
					slog.Warn("provided glob returns zero files, skipping...",
						"glob pattern", fullPath,
					)
					continue
				}
				if err != nil {
					slog.Error("failed to load templates by glob",
						"error message", err.Error(),
						"provided glob", fullPath,
					)
					continue
				}
				allTemplateFilePaths = append(allTemplateFilePaths, paths...)
			} else {
				if !utils.IsFileExists(fullPath) {
					slog.Error("failed to load template from file, file doesn't exist",
						"provided file path", fullPath)
					continue
				}
				allTemplateFilePaths = append(allTemplateFilePaths, fullPath)
			}

			if len(allTemplateFilePaths) == 0 {
				slog.Error("failed to load web ui templates, there's no files to load")
				return
			}
		}

		b.rootTempl = template.New(rootTemplateName)
		totalFiles := len(allTemplateFilePaths)
		totalLoaded := 0
		for _, f := range allTemplateFilePaths {
			err := loadTemplateFromFile(f, b.rootTempl)
			if err != nil {
				slog.Error("failed to load template file", "error message", err.Error())
				continue
			}
			totalLoaded += 1
		}

		if totalFiles != totalLoaded {
			slog.Error("not all templates were loaded in to ui bundle",
				"total files count", totalFiles,
				"total loaded", totalLoaded,
			)
		}

		logLoadedTemplates(b.rootTempl.DefinedTemplates())
	}
}

func loadTemplateFromFile(filePath string, templ *template.Template) error {
	if templ == nil {
		panic("trying to load templates for nil template")
	}
	slog.Debug(fmt.Sprintf("loading template: %s", filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	templName := utils.GetFileName(filePath, false)
	slog.Info(fmt.Sprintf("loading template %s", templName))

	_, err = templ.New(templName).Parse(string(content))
	return err
}

func logLoadedTemplates(s string) {
	if s == "" {
		return
	}
	trimmed := strings.TrimPrefix(s, "; defined templates are:")
	splited := strings.Split(trimmed, ",")
	slog.Info("successfly load templates:")
	for _, item := range splited {
		slog.Info(fmt.Sprintf("-> %s", strings.TrimSpace(item)))
	}
}
