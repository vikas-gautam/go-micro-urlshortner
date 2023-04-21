package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/data"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var FrontendDomain = "http://localhost:8080"

// // to persist the older's request data
// var mappingList = []data.URLMapping{}

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

	//get the client_ip that we have set up in frontend
	client_ip := r.Header.Get("X-Forwarded-For")

	log.Println("printing client_ip that we have set in frontend: ", client_ip)

	newMapping := data.URLMapping{Url: requestData.Url, Generated_id: id, Source_ip: client_ip}

	//save the mapping into the database
	GeneratedId, err := data.InsertUrl(newMapping)
	if err != nil {
		log.Println(err)
		return
	}

	//shortened url has been generated & saved
	var resp data.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + GeneratedId

	resp.ActualURL = requestData.Url

	//inserting this mapping in click counter
	var mappingClickCounter data.URLClickCounter

	mappingClickCounter.ShortURL = resp.ShortUrl
	mappingClickCounter.ClickCounter = 0

	//save the mapping into the database
	err = data.InsertURLClickCounter(mappingClickCounter)
	if err != nil {
		log.Println(err)
		return
	}

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
		return
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

	//shortened url has been generated & saved
	var resp data.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + id

	resp.ActualURL = actual_url

	log.Println("printing response payload", resp)

	//logic to count hits of shortURL
	var UpdateClickCounter data.URLClickCounter

	UpdateClickCounter.ShortURL = resp.ShortUrl

	err = data.UpdateURLClickCounter(UpdateClickCounter.ShortURL)
	if err != nil {
		log.Println(err)
		return
	}

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}

}

func urlClickCounter(w http.ResponseWriter, r *http.Request) {
	log.Println("resolveURL called")

	short_url := r.URL.Query().Get("u")

	//get data from the database on the basis of the id
	click_counter, err := data.GetCounterByshortUrl(short_url)
	if err != nil {
		log.Println(err)
		return
	}

	//shortened url has been generated & saved
	var resp data.ResponsePayloadUrlCounter

	// if u r using shortener-service api
	resp.ShortUrl = short_url
	resp.TotalURLClicks = click_counter

	log.Println("printing response payload", resp)

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}

}
