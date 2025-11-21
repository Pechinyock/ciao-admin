package form

import (
	"ciao-admin/cmd/coven/app/cards"
	ui "ciao-admin/cmd/coven/app/endpoints/UI"
	"ciao-admin/cmd/coven/app/projection"
	"ciao-admin/internal/utils"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const imagePoolPath = "C:/_dev/card_image_pool"

func imgPoolUploadFileFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fileGroupName := r.FormValue("group")
		if fileGroupName == "" {
			slog.Error("failed to upload image into pool grop can't be empty")
			ui.UIBundle.Render("alert", w, projection.AlertProj{
				Type:    "danger",
				Message: "Выбирите тип(группу) изображения",
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, isGroupDefinied := cards.CardTypes[fileGroupName]
		if !isGroupDefinied {
			slog.Error("failed to upload image into pool the group is not defifned",
				"request group", fileGroupName,
			)
			ui.UIBundle.Render("alert", w, projection.AlertProj{
				Type:    "danger",
				Message: "Неизвестный группа изображений",
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		overideFileName := r.FormValue("file-name")
		file, handler, _ := r.FormFile("file")
		if file == nil {
			slog.Error("failed to upload image, the file image is empty or nil")
			ui.UIBundle.Render("alert", w, projection.AlertProj{
				Type:    "danger",
				Message: "Выбирите файл",
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fileName := getFileName(handler, overideFileName)
		isPng, extension := verifyFiletype(fileName)
		if isPng {
			slog.Error("failed to upload file wrong file type, expected png", "actual", extension)
			ui.UIBundle.Render("alert", w, projection.AlertProj{
				Type:    "danger",
				Message: fmt.Sprintf("Можно загружать только png файлы, а не %s", extension),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		groupDirPath := path.Join(imagePoolPath, fileGroupName)
		if !utils.IsDirExists(groupDirPath) {
			err := os.MkdirAll(groupDirPath, 0755)
			if err != nil {
				slog.Error("failed to create dir structure", "error message", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		fullPath := path.Join(imagePoolPath, fileGroupName, fileName)
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
		ui.UIBundle.Render("alert", w, projection.AlertProj{
			Type:    "danger",
			Message: fmt.Sprintf("Не удалось загрузить файл: %s", err.Error()),
		})
		return
	}
	defer file.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		slog.Error("failed to create temporary file")
		w.WriteHeader(http.StatusInternalServerError)
		ui.UIBundle.Render("alert", w, projection.AlertProj{
			Type:    "danger",
			Message: fmt.Sprintf("Не удалось загрузить файл: %s", err.Error()),
		})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		slog.Error("failed to copy file from form data to the file")
		ui.UIBundle.Render("alert", w, projection.AlertProj{
			Type:    "danger",
			Message: fmt.Sprintf("Не удалось загрузить файл: %s", err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := utils.GetFileName(fullPath, true)
	slog.Info("the file was uploaded successfly", "file name", fileName)
	ui.UIBundle.Render("alert", w, projection.AlertProj{
		Type:    "success",
		Message: fmt.Sprintf("Файл %s успешно добавлен", fileName),
	})
}

func getFileName(handler *multipart.FileHeader, overrideFileName string) string {
	if overrideFileName == "" {
		return handler.Filename
	} else {
		extension := filepath.Ext(handler.Filename)
		result := fmt.Sprintf("%s.%s", overrideFileName, extension)
		return result
	}
}

func verifyFiletype(fileName string) (bool, string) {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	return ext != ".png", ext
}
