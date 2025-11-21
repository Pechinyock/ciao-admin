package cards

import (
	"errors"
	"log/slog"
)

/* [LAME] HARDCODE */
const CardsDataOutputPath = "C:/_dev/cards_output"

func GenerateCard(templName string, data any) error {
	if templName == "" {
		return errors.New("can't generate template with empty template name")
	}

	slog.Info("generating card", "template name", templName,
		"output path", CardsDataOutputPath,
	)
	return nil
}
