package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/data"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var Domain = "http://localhost:9090"

// to persist the older's request data
var mappingList = []data.MappingURL{}

// Healthcheck code
func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("shortener-service is healthy and ready to serve")
	var payload data.HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "shortener-service is healthy and ready to serve"

	err := writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println(err)
	}

}

func urlShortener(w http.ResponseWriter, r *http.Request) {
	log.Println("urlShortener called")
	var requestData data.RequestPayload

	//reads the request payload and save the data into requestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&requestData)
	if err != nil {
		log.Println("error while decoding", err)
	}

	//logic to shorten the actual url in the request payload
	id := uuid.New().String()[:5]

	newMapping := data.MappingURL{Url: requestData.Url, GeneratedId: id}

	//save the mapping into the database
	GeneratedId, err := data.InsertUrl(newMapping)
	if err != nil {
		log.Println(err)
		return
	}

	//shortened url has been generated & saved
	var resp data.ResponsePayload
	resp.ShortUrl = Domain + "/" + GeneratedId
	resp.Url = requestData.Url

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

	//get data from the database on the basis of the id
	actual_url, err := data.GetUrlByid(id)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, actual_url, http.StatusSeeOther)

}
