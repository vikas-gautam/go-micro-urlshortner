package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/data"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/models"

	"github.com/go-chi/chi/v5"
)

var FrontendDomain = "http://localhost:8080"

// // to persist the older's request data
// var mappingList = []models.URLMapping{}

// Healthcheck code
func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("shortener-service is healthy and ready to serve")
	var payload models.HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "shortener-service is healthy and ready to serve"

	err := writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("error while writing json response: %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
	}

}

func urlShortener(w http.ResponseWriter, r *http.Request) {
	log.Println("urlShortener called")
	var requestData models.RequestPayload

	//reads the request payload and save the data into requestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&requestData)
	if err != nil {
		log.Println("error while decoding", err)
	}

	// Initialize GeneratedId to an empty string
	GeneratedId := ""

	//get the client_ip that we have set up in frontend
	client_ip := r.Header.Get("X-Forwarded-For")
	userType := r.Header.Get("userType")

	if userType != "authenticated" {

		//Checking whether the request has already existing source ip
		counter, err := data.CheckSourceIpExistence(client_ip)
		if err == sql.ErrNoRows {
			log.Println("No rows has been matched for given source_ip (new user)\n", err)
			//that means it's new user
			GeneratedId, err = shortenURL(requestData.Url, client_ip)
			if err != nil {
				log.Println("error while generating short url for new user", err)
				msg := fmt.Sprintf("error while generating short url for new user: %v", err)
				_ = writeJSONerror(w, http.StatusInternalServerError, msg)
				return
			}
		} else if err != nil {
			log.Println("error while checking ip existence in database", err)
			msg := fmt.Sprintf("error while checking ip existence in database: %v", err)
			_ = writeJSONerror(w, http.StatusInternalServerError, msg)
			return
		}

		//check if counter is less than 5
		if counter < 5 {
			GeneratedId, err = shortenURL(requestData.Url, client_ip)
			if err != nil {
				log.Println("error while generating short url for new user", err)
				msg := fmt.Sprintf("error while generating short url for new user: %v", err)
				_ = writeJSONerror(w, http.StatusInternalServerError, msg)
				return
			}

			err = data.UpdateURLGenerateRestrictions(client_ip)
			if err != nil {
				log.Println("error while updating counter of matching source_ip", err)
				msg := fmt.Sprintf("error while updating counter of matching source_ip: %v", err)
				_ = writeJSONerror(w, http.StatusInternalServerError, msg)
				return
			}
		} else {
			err = data.UpdateStatusURLGenerateRestrictions(client_ip)
			log.Println("You have reached the limit, and status has been set to Blocked")
			if err != nil {
				log.Println("error while updating status of matching source_ip", err)
				msg := fmt.Sprintf("error while updating status of matching source_ip: %v", err)
				_ = writeJSONerror(w, http.StatusInternalServerError, msg)
				return
			}
			msg := "You have reached the limit, please signup and try again"
			_ = writeJSONerror(w, http.StatusTooManyRequests, msg)
			return
		}

		//shortened url has been generated & saved
		var resp models.ResponsePayload

		// if u r using shortener-service api
		resp.ShortUrl = FrontendDomain + "/" + GeneratedId

		resp.ActualURL = requestData.Url

		//inserting this mapping in click counter
		var mappingClickCounter models.URLClickCounter

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

	GeneratedId, err = shortenURL(requestData.Url, client_ip)
	if err != nil {
		log.Println("error while generating short url for new user", err)
		msg := fmt.Sprintf("error while generating short url for new user: %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
		return
	}
	//shortened url has been generated & saved
	var resp models.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + GeneratedId

	resp.ActualURL = requestData.Url

	//inserting this mapping in click counter
	var mappingClickCounter models.URLClickCounter

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
	if err == sql.ErrNoRows {
		log.Println("No rows has been matched for given id\n", err)
		//that means id given is not valid
		msg := "Given shortURL doesn't exist"
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return

	} else if err != nil && err != sql.ErrNoRows {
		log.Printf("error while fetching url by id: %v", err)
		msg := "Internal server error"
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
		return
	}

	//shortened url has been generated & saved
	var resp models.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + id

	resp.ActualURL = actual_url

	log.Println("printing response payload", resp)

	//logic to count hits of shortURL
	var UpdateClickCounter models.URLClickCounter

	UpdateClickCounter.ShortURL = resp.ShortUrl

	err = data.UpdateURLClickCounter(UpdateClickCounter.ShortURL)
	if err != nil {
		log.Printf("error while updating Click counter: %v", err)
		msg := fmt.Sprintln("Internal Server Error")
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
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

	if !strings.HasPrefix(short_url, "http://") && !strings.HasPrefix(short_url, "https://") {
		short_url = "http://" + short_url
	}

	//get data from the database on the basis of the id
	click_counter, err := data.GetCounterByshortUrl(short_url)
	if err == sql.ErrNoRows {
		log.Println("error while fetching Counter (shortURL): %v", err)
		//that means shortURL given doesn't exist
		msg := "Given shortURL doesn't exist"
		_ = writeJSONerror(w, http.StatusBadRequest, msg)
		return

	} else if err != nil && err != sql.ErrNoRows {
		log.Println("error while fetching Counter (shortURL): %v", err)
		msg := "Internal server error"
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
		return
	}

	//shortened url has been generated & saved
	var resp models.ResponsePayloadUrlCounter

	// if u r using shortener-service api
	resp.ShortUrl = short_url
	resp.TotalURLClicks = click_counter

	log.Println("printing response payload", resp)

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}

}
