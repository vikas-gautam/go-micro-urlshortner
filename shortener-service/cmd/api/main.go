package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/elasticsearch"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/db"
	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/storage/redis"
)

const (
	apiPort = "8080"
)

// type Config struct {
// 	DB *sql.DB
// }

func main() {

	//connect to database
	conn, err := db.ConnectToDB()
	if conn == nil {
		logrus.Panic("Can't connect to database postgres", err)
	}

	//connect to redis cache
	redisClient, err := redis.ConnectToRedis()
	if redisClient == nil && err != nil {
		logrus.Panic("Can't connect to redis", err)
	}

	//connect to elasticsearch
	esClient, err := elasticsearch.ConnectToElasticsearch()
	if err != nil {
		logrus.Panic(err)
	}

	// if u not using method
	redis.ConnectionRedis(redisClient)
	db.Connection(conn)
	elasticsearch.ConnectionElastic(esClient)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", apiPort),
		Handler: routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Can't start http server", err)
	}

	logrus.Infof("Starting shortner service on port %s\n", apiPort)

}
