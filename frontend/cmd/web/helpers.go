package main

import (
	"encoding/json"
	"net/http"

	"github.com/vikas-gautam/go-micro-urlshortner/frontend/cmd/models"
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

	var data models.AuthResponse
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
