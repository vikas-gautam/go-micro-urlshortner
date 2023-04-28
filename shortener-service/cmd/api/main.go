package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/data"
)

const (
	apiPort = "8080"
)

// type Config struct {
// 	DB *sql.DB
// }

func main() {

	//connect to database
	conn := data.ConnectToDB()
	if conn == nil {
		log.Panic("Can't connect to database postgres")
	}

	// if u not using method
	data.Connection(conn)

	// //set up config
	// app := Config{
	// 	DB: conn,
	// }

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

	log.Printf("Starting shortner service on port %s\n", apiPort)

}
