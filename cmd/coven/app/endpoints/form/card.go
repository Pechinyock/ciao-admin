package form

import (
	"ciao-admin/cmd/coven/app/cards"
	ui "ciao-admin/cmd/coven/app/endpoints/UI"
	"ciao-admin/cmd/coven/app/projection"
	"errors"
	"log/slog"
	"net/http"
	"path"
)

const cardsOutputPath = "C:/_dev/cards_output"

func cardHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createCard(w, r)
	case "PUT":
		break
	case "DELETE":
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func createCard(w http.ResponseWriter, r *http.Request) {
	selectedImageName := r.FormValue("selected-character-image")
	selectedCardType := r.FormValue("creating-card-type")
	err := distribute(selectedCardType, selectedImageName, w, r)
	if err != nil {
		slog.Error("failed to create card", "card type",
			selectedCardType, "reason", err.Error(),
		)
		ui.UIBundle.Render("alert", w, projection.AlertProj{
			Type:    "danger",
			Message: "Не удалось создать карту",
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ui.UIBundle.Render("alert", w, projection.AlertProj{
		Type:    "success",
		Message: "Карта успешно создана",
	})
}

func distribute(cardType, imageName string, w http.ResponseWriter, r *http.Request) error {
	if cardType == "" {
		return errors.New("card type is empty")
	}
	if imageName == "" {
		return errors.New("image name is empty")
	}
	switch cardType {
	case "characters":
		charName := r.FormValue("character-name")
		decorTxt := r.FormValue("character-description")
		role := r.FormValue("character-role")

		createChar(&cards.Character{
			Name:           charName,
			DecorationText: decorTxt,
			Role:           role,
			/*[LAME] HARDCODE */
			ImgPath: path.Join("image-pool", cardType, imageName),
		}, w, r)
	case "spells":
		return errors.New("unimplemented spells")
	case "secrets":
		return errors.New("unimplemented secrets")
	case "curses":
		return errors.New("unimplemented curses")
	case "ingredients":
		return errors.New("unimplemented ingredients")
	case "potions":
		return errors.New("unimplemented potions")
	default:
		return errors.New("unknown card type")
	}

	return nil
}

func createChar(cardData *cards.Character, w http.ResponseWriter, r *http.Request) {
	if cardData == nil {
		panic("trying to create card with nil data")
	}
	charsOutPath := path.Join(cardsOutputPath, "characters")
	cardData.GenerateCard(charsOutPath)
}
