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

func Index(w http.ResponseWriter, r *http.Request) {
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

func User_GET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "id must be integer")
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
		if err := writeJSONResponse(w, user); err != nil {
			panic(err)
		}
	} else {
		writeErrorResponse(w, http.StatusNotFound, "user not found")
	}
}

func User_POST(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var _json_ interface{}
	if err := json.Unmarshal(body, &_json_); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "JSON body could not be parsed")
		return
	}

	_json, ok := _json_.(map[string]interface{})
	if !ok {
		writeErrorResponse(w, http.StatusBadRequest, "JSON body must be object")
		return
	}

	db, err := OpenDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	user.CreatedAt = time.Now()
	user.UpdatedAt = nil

	if _username, ok := _json["username"]; ok {
		if username, ok := _username.(string); ok {
			user.Username = username
		} else {
			writeErrorResponse(w, http.StatusBadRequest, "'username' must be string")
			return
		}
	} else {
		writeErrorResponse(w, http.StatusBadRequest, "'username' missing")
		return
	}

	if _email, ok := _json["email"]; ok {
		if email, ok := _email.(string); ok {
			user.Email = email
		} else {
			writeErrorResponse(w, http.StatusBadRequest, "'email' must be string")
			return
		}
	} else {
		writeErrorResponse(w, http.StatusBadRequest, "'email' missing")
		return
	}

	var userRole UserRole
	if _roleLabel, ok := _json["roleLabel"]; ok {
		if roleLabel, ok := _roleLabel.(string); ok {
			db.Where(UserRole{Label: roleLabel}).Take(&userRole)

			if userRole.ID <= 0 {
				writeErrorResponse(w, http.StatusNotFound, "role not found")
				return
			}
		} else {
			writeErrorResponse(w, http.StatusBadRequest, "'roleLabel' must be string")
			return
		}
	} else {
		writeErrorResponse(w, http.StatusBadRequest, "'roleLabel' missing")
		return
	}
	user.RoleId = userRole.ID

	var userAddress UserAddress
	if _address, ok := _json["address"]; ok {
		if address, ok := _address.(map[string]interface{}); ok {
			if _iAddress, ok := address["address"]; ok {
				if _iAddress == nil {
					userAddress.Address = nil
				} else if iAddress, ok := _iAddress.(string); ok {
					userAddress.Address = &iAddress
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "inner 'address' must be string or null")
					return
				}
			}

			if _province, ok := address["province"]; ok {
				if _province == nil {
					userAddress.Province = nil
				} else if province, ok := _province.(string); ok {
					userAddress.Province = &province
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "'province' must be string or null")
					return
				}
			}
		} else {
			writeErrorResponse(w, http.StatusBadRequest, "'address' must be object")
			return
		}
	} else {
		writeErrorResponse(w, http.StatusBadRequest, "'address' missing")
		return
	}
	user.Address = userAddress

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

	writeSuccessResponse(w)
}

func User_PUT(w http.ResponseWriter, r *http.Request) {

}

func User_DELETE(w http.ResponseWriter, r *http.Request) {

}

func Users_GET(w http.ResponseWriter, r *http.Request) {

}
