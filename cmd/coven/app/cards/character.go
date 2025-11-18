package cards

import "log/slog"

const templatePath = "C:/_dev/card_templates"

type Character struct {
	Name           string `json:"name"`
	Role           string `json:"role"`
	DecorationText string `json:"decorationText"`
	ImgPath        string `json:"imagePath"`
}

func (c *Character) GenerateCard(outPath string) error {
	slog.Info("generating card", "template path", templatePath)
	return nil
}
