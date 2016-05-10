package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	. "github.com/danielleknudson/user-base/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func CreateUser(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	user := &User{
		FirstName: req.FormValue("first_name"),
		LastName:  req.FormValue("last_name"),
		Email:     req.FormValue("email"),
	}

	password, err := HashPassword(req.FormValue("password"))

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	user.Password = password
	err = user.Save()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(user)
}

func FetchUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	userId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user := new(User)
	user.Id = userId

	err = user.Fetch()

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(struct {
		*User
		Password bool `json:"password,omitempty"`
	}{
		User: user,
	})
}

func FetchUsers(w http.ResponseWriter, req *http.Request) {
	users := NewUsers()
	publicUsers, err := users.FetchAll()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	json.NewEncoder(w).Encode(publicUsers)
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	user := NewUser()
	vars := mux.Vars(req)
	id := vars["id"]
	userId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user.Id = userId

	// Fetch the user's existing data
	user.Fetch()

	err = req.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	decoder := schema.NewDecoder()

	// Add the user's new data to the struct
	err = decoder.Decode(user, req.PostForm)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Save the user's changes
	err = user.Update()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DestroyUser(w http.ResponseWriter, req *http.Request) {
	user := NewUser()
	vars := mux.Vars(req)
	id := vars["id"]

	userId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user.Id = userId

	err = user.Destroy()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
