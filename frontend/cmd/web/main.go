package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	webPort = "8080"
)

func main() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Starting urlshortener frontend  on port %s\n", webPort)
}
