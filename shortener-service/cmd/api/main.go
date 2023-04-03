package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/data"
)

const (
	apiPort = "8080"
)

var counts int64

// type Config struct {
// 	DB     *sql.DB
// 	Models data.Models
// }

func main() {

	//connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database postgres")
	}

	data.Connection(conn)

	// //set up config
	// app := Config{
	// 	DB:     conn,
	// 	Models: data.New(conn),
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet...", err)
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue

	}
}
