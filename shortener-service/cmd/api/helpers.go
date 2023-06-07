package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/models"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/db"
)

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func writeJSONerror(w http.ResponseWriter, status int, msg string) error {

	var data models.ResponsePayloadError
	data.Status = status
	data.Message = msg

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// actual logic to shorten the url
func shortenURL(request_url, client_ip, user_type, email string) (string, error) {

	//logic to shorten the actual url in the request payload
	id := uuid.New().String()[:5]

	log.Println("printing client_ip that we have set in frontend: ", client_ip)

	newMapping := models.URLMapping{Url: request_url, Generated_id: id, Source_ip: client_ip, User_type: user_type, Email: email}

	//save the mapping into the database
	GeneratedId, err := db.InsertUrl(newMapping)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return GeneratedId, nil

}
