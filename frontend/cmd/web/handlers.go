package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/vikas-gautam/go-micro-urlshortner/frontend/cmd/models"
)

//go:embed templates
var templateFS embed.FS

// Healthcheck code
func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving urlshortener homepage")

	render(w, "index.html")

}

// urlshortener test
func urlShortener(w http.ResponseWriter, r *http.Request) {

	//checking whether the user is guest or signedup
	username, password, ok := r.BasicAuth()
	log.Printf("username: %v, password: %v, ok: %v", username, password, ok)

	//fetch client_ip from request
	client_ip := r.Header.Get("X-Forwarded-For")

	//creating backend url
	BACKEND_SERVICE := os.Getenv("BACKEND_API_URL") + "/getshortenurl"

	req, _ := http.NewRequest("POST", BACKEND_SERVICE, r.Body)

	if ok {
		resp := userauth(username, password)

		if resp.Status != 200 {
			_ = writeJSON(w, resp.Status, resp)
			return
		} else {
			fmt.Fprintf(w, "Authenticated successfully, welcome %s\n", username)
			// Set the X-Forwarded-For header
			req.Header.Set("X-Forwarded-For", client_ip)
			req.Header.Set("userType", "authenticated")
		}

	} else {
		fmt.Printf("Basic Auth header not present so proceeding as a guest user")
		// Set the X-Forwarded-For header
		req.Header.Set("X-Forwarded-For", client_ip)
		req.Header.Set("userType", "guest")

	}

	log.Println("fetching response from urlshortener service")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error while making request to backend shortener svc", err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	if response.StatusCode != 200 {
		var resp models.AuthResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			panic(err)
		}
		_ = writeJSONerror(w, resp.Status, resp.Message)
		return

	} else {

		// Unmarshal JSON response into URLCollection struct
		var resp models.URLCollection

		err = json.Unmarshal(body, &resp)
		if err != nil {
			panic(err)
		}
		log.Println(resp)

		successResp := models.SuccessResponse{
			Response: resp,
		}
		err = writeJSON(w, http.StatusOK, successResp)
		if err != nil {
			log.Println("error while writing json response to frontend request", err)
		}

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

	//creating backend url
	BACKEND_SERVICE := os.Getenv("BACKEND_API_URL") + "/signup"

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
	var resp models.URLCollection
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	log.Println("printing actual url fetched from backend - ", resp.ActualURL)

	rediected_url := resp.ActualURL
	if !strings.HasPrefix(rediected_url, "http://") && !strings.HasPrefix(rediected_url, "https://") {
		rediected_url = "http://" + rediected_url
	}

	http.Redirect(w, r, rediected_url, http.StatusSeeOther)

}

func signup(w http.ResponseWriter, r *http.Request) {
	log.Println("fetching response from auth service to get authenticated")

	//creating backend url
	AUTH_SERVICE := os.Getenv("AUTH_API_URL") + "/user/signup"

	req, _ := http.NewRequest("POST", AUTH_SERVICE, r.Body)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error while making request to auth svc", err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// Unmarshal JSON response into  struct
	var resp models.AuthResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	err = writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println("error while writing json response to frontend request", err)
	}

}

func userauth(username, password string) models.AuthResponse {

	//creating backend url
	AUTH_SERVICE := os.Getenv("AUTH_API_URL") + "/user/auth"

	req, err := http.NewRequest("GET", AUTH_SERVICE, nil)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	u := username
	p := password
	auth := u + ":" + p
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", basicAuth)

	log.Println("hitting auth service to login")

	log.Println(AUTH_SERVICE)
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	// Unmarshal JSON response into auth struct
	var resp models.AuthResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	log.Println(resp)
	return resp

}
