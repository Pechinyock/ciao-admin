package form

import (
	"ciao-admin/internal/utils"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const imagePoolPath = "C:/_dev/card_image_pool"
const characterImagesDirName = "characters"

func imgPoolUploadFileFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fileName := getFileName(r)
		isPng, extension := verifyFiletype(fileName)
		if isPng {
			slog.Error("failed to upload file wrong file type, expected png", "actual", extension)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fullPath := path.Join(imagePoolPath, characterImagesDirName, fileName)
		uploadImage(fullPath, w, r)
	case "GET":
	case "DELETE":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func uploadImage(fullPath string, w http.ResponseWriter, r *http.Request) {
	/* Maximum 10 mb */
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		slog.Error("failed to upload file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		slog.Error("failed to create temporary file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		slog.Error("failed to copy file from form data to the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := utils.GetFileName(fullPath, true)
	slog.Info("the file was uploaded successfly", "file name", fileName)
}

func getFileName(r *http.Request) string {
	overideFileName := r.FormValue("file-name")
	file, handler, _ := r.FormFile("file")
	if file != nil {

	}
	if overideFileName == "" {
		return handler.Filename
	} else {
		extension := filepath.Ext(handler.Filename)
		result := fmt.Sprintf("%s.%s", overideFileName, extension)
		return result
	}
}

func verifyFiletype(fileName string) (bool, string) {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	return ext != ".png", ext
}
