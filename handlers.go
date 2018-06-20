package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func User_GET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := OpenDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	db.Preload("Role").Preload("Address").Table("users").First(&user, id)

	if user.ID == id {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func User_POST(w http.ResponseWriter, r *http.Request) {

}

func User_PUT(w http.ResponseWriter, r *http.Request) {

}

func User_DELETE(w http.ResponseWriter, r *http.Request) {

}

func Users_GET(w http.ResponseWriter, r *http.Request) {

}
