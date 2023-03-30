package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var Domain = "http://localhost:8080"

// Healthcheck code
type HealthPayload struct {
	Message string
	Status  int
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("shortner-service is healthy and ready to serve")
	var payload HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "shortner-service is healthy and ready to serve"

	//w.Write requires []byte data
	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}

}

// To take the user's Post request #############################################

type RequestPayload struct {
	Url string `json:"url"`
}

type ResponsePayload struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

type MappingURL struct {
	Url         string `json:"url"`
	GeneratedId string `json:"generatedId"`
}

// to persist the older's request data
var mappingList = []MappingURL{}

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

	//w.Write requires []byte data
	out, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Println("error while marshal", err)
	}

	//returns the shortened url to user as response payload
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}

}

//user hits shortenedURL it should redirect to actual url #########################
// Logic to resolve generated shortenedURL ########################################

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
