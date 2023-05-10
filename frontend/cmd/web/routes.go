package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {

	mux := chi.NewRouter()
	//specifiy who is allowed to connect to
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", homepage)
	mux.Post("/getShortUrl", urlShortener)
	mux.Get("/{id}", resolveURL)

	mux.Post("/user/signup", signup)
	mux.Get("/user/login", login)

	return mux

}
