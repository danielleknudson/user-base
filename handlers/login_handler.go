package handlers

import (
	"net/http"

	. "github.com/danielleknudson/user-base/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := &User{
		Email: r.FormValue("email"),
	}

	err := user.FindByEmail()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	password := r.FormValue("password")

	validPassword := ComparePasswords(user.Password, password)

	if !validPassword {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	w.WriteHeader(http.StatusOK)
}
