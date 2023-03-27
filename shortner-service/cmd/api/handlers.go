package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
