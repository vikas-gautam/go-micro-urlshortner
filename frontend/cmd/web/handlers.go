package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

// Healthcheck code
func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving urlshortener homepage")

	render(w, "index.html")

}

//go:embed templates
var templateFS embed.FS

type URLCollection struct {
	ActualURL string
	ShortURL  string
}
type SuccessResponse struct {
	Response URLCollection
}

// urlshortener test
func urlShortener(w http.ResponseWriter, r *http.Request) {
	log.Println("fetching response from urlshortener service")

	req, _ := http.NewRequest("POST", "http://shortener-service:8080/getshortenurl", r.Body)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error while making request to backend shortener svc", err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// Unmarshal JSON response into URLCollection struct
	var resp URLCollection
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	log.Println(resp)

	successResp := SuccessResponse{
		Response: resp,
	}

	err = writeJSON(w, http.StatusOK, successResp)
	if err != nil {
		log.Println("error while writing json response to frontend request", err)
	}

}

func render(w http.ResponseWriter, t string) {

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func resolveURL(w http.ResponseWriter, r *http.Request) {

	log.Println("fetching response from urlshortener service to resolve shorten url")

	id := chi.URLParam(r, "id")
	log.Println("printing random number id of shortenedURL=>", id)

	BACKEND_SERVICE := "http://shortener-service:8080" + fmt.Sprintf("/"+id)

	log.Println(BACKEND_SERVICE)

	response, err := http.Get(BACKEND_SERVICE)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// Unmarshal JSON response into URLCollection struct
	var resp URLCollection
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	log.Println("priting actual url fetched from backend - ", resp.ActualURL)

	http.Redirect(w, r, resp.ActualURL, http.StatusSeeOther)

}