package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/models"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/db"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/redis"
)

var FrontendDomain = "http://localhost:8080"

// Healthcheck code
func healthCheck(w http.ResponseWriter, r *http.Request) {
	var payload models.HealthPayload

	payload.Status = http.StatusOK
	payload.Message = "shortener-service is healthy and ready to serve"

	err := writeJSON(w, http.StatusOK, payload)
	if err != nil {
		logrus.Error("error while writing json response:", err)
		msg := fmt.Sprintf("error while writing json response: %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, msg)
	}

}

func urlShortener(w http.ResponseWriter, r *http.Request) {

	var requestData models.RequestPayload

	//reads the request payload and save the data into requestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&requestData)
	if err != nil {
		logrus.Error("error while decoding", err)
	}

	//get the client_ip that we have set up in frontend
	client_ip := r.Header.Get("X-Forwarded-For")
	userType := r.Header.Get("userType")
	email := r.Header.Get("email")

	if userType == "Guest" {
		counter, err := db.CheckSourceIpExistence(client_ip)
		logrus.Info("Printing counter while checking source_ip existence:  ", counter)

		if err != sql.ErrNoRows && err != nil {
			_ = writeJSONerror(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		if counter >= 5 {
			err = db.UpdateStatusURLGenerateRestrictions(client_ip)
			logrus.Error("Error while calling UpdateStatusURLGenerateRestrictions", err)
			_ = writeJSONerror(w, http.StatusTooManyRequests, "You have reached the limit, please signup and try again")
			return
		}

		//this is handling new user i.e. err == no rows
		err = db.UpdateURLGenerateRestrictions(client_ip)
		if err != nil {
			logrus.Error(err)
		}
	}

	GeneratedId, _ := shortenURL(requestData.Url, client_ip, userType, email)

	//this will run when the short url has been generated

	//shortened url has been generated & saved
	var resp models.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + GeneratedId

	resp.ActualURL = requestData.Url

	//inserting the same in redis
	err = redis.SetData(GeneratedId, resp.ActualURL)
	if err != nil {
		logrus.Error(err)
		return
	}

	//inserting this mapping in click counter
	var mappingClickCounter models.URLClickCounter
	mappingClickCounter.ShortURL = resp.ShortUrl
	mappingClickCounter.ClickCounter = 0

	//save the mapping into the database
	err = db.InsertURLClickCounter(mappingClickCounter)
	if err != nil {
		logrus.Error(err)
		return
	}

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		logrus.Error(err)
		return
	}

}

//user hits shortenedURL it should redirect to actual url
// Logic to resolve generated shortenedURL

func resolveURL(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// Fetch data from Redis
	actual_url, err := redis.GetData(id)

	// Check if actual_url is nil and fetch from the database if needed
	if actual_url == "" || err != nil {
		// Data not found in Redis, fetch from the database
		actual_url, err = db.GetUrlByid(id)

		// Handle the error if any
		if err == sql.ErrNoRows {
			// That means the given id is not valid
			logrus.Error("Error while fetching URL by id: ", err)
			_ = writeJSONerror(w, http.StatusBadRequest, "Given shortURL doesn't exist")
			return
		} else if err != nil {
			logrus.Error("Error while fetching URL by id: ", err)
			_ = writeJSONerror(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		//inserting the same in redis
		err = redis.SetData(id, actual_url)
		if err != nil {
			log.Println(err)
			return
		}
		logrus.Info("Cache MISS")
	} else {
		logrus.Info("Cache HIT")
	}

	//shortened url has been generated & saved
	var resp models.ResponsePayload

	// if u r using shortener-service api
	resp.ShortUrl = FrontendDomain + "/" + id

	resp.ActualURL = actual_url

	//logic to count hits of shortURL
	var UpdateClickCounter models.URLClickCounter

	UpdateClickCounter.ShortURL = resp.ShortUrl

	err = db.UpdateURLClickCounter(UpdateClickCounter.ShortURL)
	if err != nil {

		logrus.Error("error while updating Click counter: ", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		logrus.Error(err)
	}

}

func urlClickCounter(w http.ResponseWriter, r *http.Request) {

	short_url := r.URL.Query().Get("u")

	if !strings.HasPrefix(short_url, "http://") && !strings.HasPrefix(short_url, "https://") {
		short_url = "http://" + short_url
	}

	//get data from the database on the basis of the id
	click_counter, err := db.GetCounterByshortUrl(short_url)
	if err == sql.ErrNoRows {
		logrus.Errorf("error while fetching Counter (shortURL): %v", err)
		//that means shortURL given doesn't exist
		_ = writeJSONerror(w, http.StatusBadRequest, "Given shortURL doesn't exist")
		return

	} else if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("error while fetching Counter (shortURL): %v", err)
		_ = writeJSONerror(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	//shortened url has been generated & saved
	var resp models.ResponsePayloadUrlCounter

	// if u r using shortener-service api
	resp.ShortUrl = short_url
	resp.TotalURLClicks = click_counter

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}

}
