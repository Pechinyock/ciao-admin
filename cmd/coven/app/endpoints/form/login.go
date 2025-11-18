package form

import (
	"ciao-admin/cmd/coven/app/repository"
	"ciao-admin/cmd/coven/app/security"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func loginHandleFunc(w http.ResponseWriter, r *http.Request) {
	checkPasswordFunc := func(password, hash string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			slog.Error(err.Error())
		}
		return err == nil
	}
	login := r.PostFormValue("login")
	pwd := r.PostFormValue("password")

	usr, err := repository.UserRepo.Get(login)
	if err != nil {
		slog.Error(err.Error())
	}

	if !checkPasswordFunc(pwd, usr.PasswordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	JwtProvider, err := security.NewTokenProvider()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := JwtProvider.GenerateToken(usr.Login)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setJwtCookie(token, w)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func setJwtCookie(jwtStr string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Domain:   "localhost",
		Path:     "/",
		Name:     "coven-token",
		Value:    jwtStr,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})
}
