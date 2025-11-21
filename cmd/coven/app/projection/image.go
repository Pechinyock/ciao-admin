package projection

type ImageViewProj struct {
	BasePath  string
	FileGroup string
	Images    []ImageProj
}

type ImageProj struct {
	FileName string
}
