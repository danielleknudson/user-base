package main

import (
	"log"
	"net/http"

	. "github.com/danielleknudson/user-base/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", FetchUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", FetchUser).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", DestroyUser).Methods("DELETE")

	r.HandleFunc("/signup", CreateUser).Methods("POST")

	r.HandleFunc("/login", LoginHandler).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
