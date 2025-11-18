package form

import (
	"ciao-admin/cmd/coven/app/cards"
	"net/http"
)

const outputPath = "C:/_dev/cards_output"

func characterHadleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createChar(w, r)
	case "PUT":
		break
	case "DELETE":
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func createChar(w http.ResponseWriter, r *http.Request) {
	charName := r.PostFormValue("character-name")
	charRole := r.PostFormValue("character-role")
	decorationText := r.PostFormValue("character-description")
	imagePath := r.PostFormValue("character-img-path")
	char := cards.Character{
		Name:           charName,
		Role:           charRole,
		DecorationText: decorationText,
		ImgPath:        imagePath,
	}
	char.GenerateCard(outputPath)
}
