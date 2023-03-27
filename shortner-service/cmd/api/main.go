package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	apiPort = "8080"
)

func main() {
	log.Printf("Starting shortner service on port %s\n", apiPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", apiPort),
		Handler: routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
