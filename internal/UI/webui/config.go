package webui

type WebUIBundleConfig struct {
	RootPath                 string   `json:"rootPath"`
	StaticFilesDirName       string   `json:"staticFilesDirName"`
	StaticFilesRootRouteName string   `json:"staticFilesRootRouteName"`
	TemplatePaths            []string `json:"templatePaths"`
}
