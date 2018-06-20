package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	maxBodyBytes = 1048576
)

func getIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func writeJSONResponse(w http.ResponseWriter, v interface{}) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func writeSuccessResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	io.WriteString(w, message)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "id must be integer")
		return
	}

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	db.Preload("Role").Preload("Address").Table("users").First(&user, id)

	if user.ID == id {
		if err := writeJSONResponse(w, user); err != nil {
			panic(err)
		}
	} else {
		writeErrorResponse(w, http.StatusNotFound, "user not found")
	}
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var jsonMap map[string]interface{}
	{
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			writeErrorResponse(w, http.StatusUnprocessableEntity, "JSON body could not be parsed")
			return
		}

		var ok bool
		jsonMap, ok = jsonBody.(map[string]interface{})
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "JSON body must be object")
			return
		}
	}

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	user.CreatedAt = time.Now()
	user.UpdatedAt = nil

	{
		tUsername, ok := jsonMap["username"]
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'username' missing")
			return
		}

		username, ok := tUsername.(string)
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'username' must be string")
			return
		}

		user.Username = username
	}

	{
		tEmail, ok := jsonMap["email"]
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'email' missing")
			return
		}

		email, ok := tEmail.(string)
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'email' must be string")
			return
		}

		user.Email = email
	}

	{
		var userRole UserRole

		tRoleLabel, ok := jsonMap["roleLabel"]
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'roleLabel' missing")
			return
		}

		roleLabel, ok := tRoleLabel.(string)
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'roleLabel' must be string")
			return
		}

		db.Where(UserRole{Label: roleLabel}).Take(&userRole)

		if userRole.ID <= 0 {
			writeErrorResponse(w, http.StatusNotFound, "role not found")
			return
		}

		user.RoleID = userRole.ID
	}

	{
		var userAddress UserAddress

		tAddressMap, ok := jsonMap["address"]
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'address' missing")
			return
		}

		addressMap, ok := tAddressMap.(map[string]interface{})
		if !ok {
			writeErrorResponse(w, http.StatusBadRequest, "'address' must be object")
			return
		}

		tAddress, ok := addressMap["address"]
		if ok {
			if tAddress == nil {
				userAddress.Address = nil
			} else if iAddress, ok := tAddress.(string); ok {
				userAddress.Address = &iAddress
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "inner 'address' must be string or null")
				return
			}
		}

		tProvince, ok := addressMap["province"]
		if ok {
			if tProvince == nil {
				userAddress.Province = nil
			} else if province, ok := tProvince.(string); ok {
				userAddress.Province = &province
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'province' must be string or null")
				return
			}
		}
		user.Address = userAddress
	}

	{
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if tx.Error != nil {
			panic(tx.Error)
		}

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			writeErrorResponse(w, http.StatusConflict, "user already exists")
			return
		}

		if err := tx.Commit().Error; err != nil {
			panic(err)
		}
	}

	writeSuccessResponse(w)
}

func putUserHander(w http.ResponseWriter, r *http.Request) {

}

func deleteUserHander(w http.ResponseWriter, r *http.Request) {

}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

}
