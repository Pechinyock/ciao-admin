package ui

import (
	"ciao-admin/cmd/coven/app/projection"
	"errors"
	"os"
	"path"
	"strings"
)

/*[LAME] HARDCODE */
const physicalBasePath = "C:/_dev/card_image_pool"
const baseUriPath = "image-pool"

func loadImagesPrewiewData(poolGroupName string) ([]projection.ImageProj, error) {
	if poolGroupName == "" {
		return nil, errors.New("pool group name is empty")
	}
	pathToDir := path.Join(physicalBasePath, poolGroupName)
	files, err := os.ReadDir(pathToDir)
	if err != nil {
		return nil, err
	}
	var result []projection.ImageProj
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
			added := projection.ImageProj{
				FileName: file.Name(),
			}
			result = append(result, added)
		}
	}
	return result, nil
}
