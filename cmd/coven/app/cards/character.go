package cards

import (
	"ciao-admin/internal/utils"
	"errors"
	"html/template"
	"log/slog"
	"os"
	"path"
	"strings"
)

/* [LAME] HARDCODE */
const templatePath = "C:/_dev/card_templates"
const baseUriPath = "image-pool"

type Character struct {
	Name           string `json:"name"`
	Role           string `json:"role"`
	DecorationText string `json:"decorationText"`
	ImgPath        string `json:"imagePath"`
}

func (c *Character) GenerateCard(outPath string) error {
	if c.Name == "" {
		return errors.New("character name couldn't be empty string")
	}
	slog.Info("generating card", "template path", templatePath)
	characterTemplatePath := path.Join(templatePath, "character_templ.html")
	if !utils.IsFileExists(characterTemplatePath) {
		return errors.New("character template file deosn't exitst")
	}
	if c.ImgPath != "" {
		c.ImgPath = getBackgroudSource(c.ImgPath)
	}
	templateName := "character-card"
	templ, err := template.New(templateName).ParseFiles(characterTemplatePath)
	if err != nil {
		return err
	}
	if utils.IsFilePath(outPath) {
		return errors.New("provided out path is a path to the file")
	}
	if !utils.IsDirExists(outPath) {
		if err := os.MkdirAll(outPath, 0755); err != nil {
			return err
		}
	}

	fullPath := path.Join(outPath, c.Name)
	if utils.IsFileExists(fullPath) {
		slog.Warn("overriding existg card", "path", fullPath)
	}
	fileResult, err := os.Create(fullPath + ".html")
	if err != nil {
		return err
	}
	defer fileResult.Close()

	err = templ.ExecuteTemplate(fileResult, templateName, c)
	if err != nil {
		return err
	}
	return nil
}

/*[LAME]*/
func getBackgroudSource(incomingImgPath string) string {
	remoteSrc := "image-pool"
	if strings.HasSuffix(remoteSrc, remoteSrc) {
		return path.Join("https://localhost:6969", incomingImgPath)
	} else {
		return path.Join("C:/_dev/card_image_pool")
	}
}
