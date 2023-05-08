package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vikas-gautam/go-micro-urlshortner/auth-service/cmd/data"
	"github.com/vikas-gautam/go-micro-urlshortner/auth-service/cmd/models"
)

var FrontendDomain = "http://localhost:8080"

var Validate = validator.New()

// Healthcheck code
func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("auth-service is healthy and ready to serve")
	var payload models.HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "auth-service is healthy and ready to serve"

	err := writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("error while writing json response: %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
	}

}

// login function
func login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("Error parsing basic auth")
		msg := "error parsing basic auth"
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return
	}
	hashedPassword, _ := data.GetHashedPasswordByEmailid(username)
	match := CheckPasswordHash(password, hashedPassword)
	if match {
		msg := "you are successfully logged in"
		_ = writeJSONerror(w, http.StatusOK, msg)
	}
}

// signup func
func signup(w http.ResponseWriter, r *http.Request) {
	var signupPayload models.Users

	//reads the request payload and save the data into signupPayload
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&signupPayload)
	if err != nil {
		log.Println("error while decoding", err)
		msg := fmt.Sprintf("error while reading json request: %v", err)
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	validationErr := Validate.Struct(signupPayload)
	if validationErr != nil {
		log.Println("error while validating", validationErr)
		msg := fmt.Sprintf("error while validating fields in request: %v", validationErr)
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return
	}

	//logic to check if user already exists or not
	emailExists, err := data.GetUserByEmailid(signupPayload.Email)
	if err != nil && err == sql.ErrNoRows {
		log.Println(err)
	}

	//logic to check if user already exists or not
	phoneExists, err := data.GetUserByPhone(signupPayload.Phone)
	if err != nil && err == sql.ErrNoRows {
		log.Println(err)
	}

	//email not existing and phone exists
	if !emailExists && !phoneExists {
		//insert data into the database
		signupPayload.Status = "Active"

		//make hash out of given password to avoid saving actual password
		hashedPassword, _ := HashPassword(signupPayload.Password)
		signupPayload.Password = hashedPassword

		err = data.InsertSignupData(signupPayload)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		msg := ""
		if emailExists {
			msg = fmt.Sprintf("email id: %s already exists", signupPayload.Email)
		}
		if phoneExists {
			if msg != "" {
				msg += " and "
			}
			msg += fmt.Sprintf("phone number: %s already exists", signupPayload.Phone)
		}
		log.Println(msg)
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return
	}

	resp := fmt.Sprintf("You are successfully signedup and here is your Username: %s", signupPayload.Email)
	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("error while writing json response: %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
		return
	}

}
