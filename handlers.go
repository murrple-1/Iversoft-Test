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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	db.Preload("Role").Preload("Address").Take(&user, id)

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

		if tAddress, ok := addressMap["address"]; ok {
			if tAddress == nil {
				userAddress.Address = nil
			} else if iAddress, ok := tAddress.(string); ok {
				userAddress.Address = &iAddress
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "inner 'address' must be string or null")
				return
			}
		}

		if tProvince, ok := addressMap["province"]; ok {
			if tProvince == nil {
				userAddress.Province = nil
			} else if province, ok := tProvince.(string); ok {
				userAddress.Province = &province
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'province' must be string or null")
				return
			}
		}

		if tCity, ok := addressMap["city"]; ok {
			if tCity == nil {
				userAddress.City = nil
			} else if city, ok := tCity.(string); ok {
				userAddress.City = &city
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'city' must be string or null")
				return
			}
		}

		if tCountry, ok := addressMap["country"]; ok {
			if tCountry == nil {
				userAddress.Country = nil
			} else if country, ok := tCountry.(string); ok {
				userAddress.Country = &country
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'country' must be string or null")
				return
			}
		}

		if tPostalCode, ok := addressMap["postalCode"]; ok {
			if tPostalCode == nil {
				userAddress.PostalCode = nil
			} else if postalCode, ok := tPostalCode.(string); ok {
				userAddress.PostalCode = &postalCode
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'postalCode' must be string or null")
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
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "id must be integer")
		return
	}

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
	db.Preload("Address").Take(&user, id)

	hasChanged := false

	{
		if tEmail, ok := jsonMap["email"]; ok {
			if email, ok := tEmail.(string); ok {
				user.Email = email

				hasChanged = true
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'email' must be string")
				return
			}
		}
	}

	{
		if tRoleLabel, ok := jsonMap["roleLabel"]; ok {
			if roleLabel, ok := tRoleLabel.(string); !ok {
				var userRole UserRole
				db.Where(UserRole{Label: roleLabel}).Take(&userRole)

				if userRole.ID >= 1 {
					user.RoleID = userRole.ID

					hasChanged = true
				} else {
					writeErrorResponse(w, http.StatusNotFound, "user role not found")
					return
				}
			} else {
				writeErrorResponse(w, http.StatusBadRequest, "'roleLabel' must be string")
				return
			}
		}
	}

	{
		if tAddressMap, ok := jsonMap["address"]; ok {
			addressMap, ok := tAddressMap.(map[string]interface{})
			if !ok {
				writeErrorResponse(w, http.StatusBadRequest, "'address' must be object")
				return
			}

			if tAddress, ok := addressMap["address"]; ok {
				if tAddress == nil {
					user.Address.Address = nil
				} else if address, ok := tAddress.(string); ok {
					user.Address.Address = &address
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "inner 'address' must be string or null")
					return
				}

				hasChanged = true
			}

			if tProvince, ok := addressMap["province"]; ok {
				if tProvince == nil {
					user.Address.Province = nil
				} else if province, ok := tProvince.(string); ok {
					user.Address.Province = &province
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "'province' must be string or null")
					return
				}

				hasChanged = true
			}

			if tCity, ok := addressMap["city"]; ok {
				if tCity == nil {
					user.Address.City = nil
				} else if city, ok := tCity.(string); ok {
					user.Address.City = &city
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "'city' must be string or null")
					return
				}

				hasChanged = true
			}

			if tCountry, ok := addressMap["country"]; ok {
				if tCountry == nil {
					user.Address.Country = nil
				} else if country, ok := tCountry.(string); ok {
					user.Address.Country = &country
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "'country' must be string or null")
					return
				}

				hasChanged = true
			}

			if tPostalCode, ok := addressMap["postalCode"]; ok {
				if tPostalCode == nil {
					user.Address.PostalCode = nil
				} else if postalCode, ok := tPostalCode.(string); ok {
					user.Address.PostalCode = &postalCode
				} else {
					writeErrorResponse(w, http.StatusBadRequest, "'postalCode' must be string or null")
					return
				}

				hasChanged = true
			}
		}
	}

	if hasChanged {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if tx.Error != nil {
			panic(tx.Error)
		}

		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			writeErrorResponse(w, http.StatusConflict, "user conflict")
			return
		}

		if err := tx.Commit().Error; err != nil {
			panic(err)
		}
	}

	writeSuccessResponse(w)
}

func deleteUserHander(w http.ResponseWriter, r *http.Request) {
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
	db.Take(&user, id)

	if user.ID <= 0 {
		writeErrorResponse(w, http.StatusNotFound, "user not found")
		return
	}

	address := user.Address

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

		if err := tx.Delete(&user).Error; err != nil {
			panic(err)
		}

		if err := tx.Delete(&address).Error; err != nil {
			panic(err)
		}

		if err := tx.Commit().Error; err != nil {
			panic(err)
		}
	}

	writeSuccessResponse(w)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var users []User
	db.Preload("Role").Preload("Address").Find(&users)

	if err := writeJSONResponse(w, users); err != nil {
		panic(err)
	}
}
