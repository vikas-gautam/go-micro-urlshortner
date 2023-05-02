package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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

	err = writeJSON(w, http.StatusOK, signupPayload)
	if err != nil {
		log.Println(err)
		return
	}

}
