package webui

type WebUIBundleConfig struct {
	RootPath                 string   `json:"rootPath"`
	StaticFilesDirName       string   `json:"staticFilesDirName"`
	StaticFilesRootRouteName string   `json:"staticFilesRootRouteName"`
	LayoutTemplates          []string `json:"layoutTemplates"`
	PageTemplates            []string `json:"pageTemplates"`
}
