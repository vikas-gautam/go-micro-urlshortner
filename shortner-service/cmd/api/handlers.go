package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var Domain = "http://localhost:8080"

// to persist the older's request data
var mappingList = []MappingURL{}

// Healthcheck code
func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("shortner-service is healthy and ready to serve")
	var payload HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "shortner-service is healthy and ready to serve"

	err := writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println(err)
	}

}

func urlShortener(w http.ResponseWriter, r *http.Request) {
	log.Println("urlShortener called")
	var requestData RequestPayload

	//reads the request payload and save the data into requestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&requestData)
	if err != nil {
		log.Println("error while decoding", err)
	}

	//logic to shorten the actual url in the request payload
	id := uuid.New().String()[:5]

	var resp ResponsePayload
	resp.Url = requestData.Url
	resp.ShortUrl = Domain + "/" + id

	newMapping := MappingURL{Url: requestData.Url, GeneratedId: id}

	mappingList = append(mappingList, newMapping)

	log.Println("printing whole  slice", mappingList)

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}

}

//user hits shortenedURL it should redirect to actual url
// Logic to resolve generated shortenedURL

func resolveURL(w http.ResponseWriter, r *http.Request) {
	log.Println("resolveURL called")

	id := chi.URLParam(r, "id")
	log.Println("printing random number id of shortenedURL=>", id)

	foundMatch := false
	for _, mappedURL := range mappingList {
		if mappedURL.GeneratedId == id {
			foundMatch = true
			log.Println("Redirecting to actual URL", mappedURL.Url)
			http.Redirect(w, r, mappedURL.Url, http.StatusSeeOther)
			break
		}
	}
	if !foundMatch {
		log.Println("Please provide a valid shortenedURL")
	}

}
