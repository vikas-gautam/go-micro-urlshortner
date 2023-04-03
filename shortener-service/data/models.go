package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type MappingURLFromDB struct {
	Url         string `json:"url"`
	GeneratedId string `json:"generatedId"`
}

type HealthPayload struct {
	Message string
	Status  int
}

// To take the user's Post request #############################################
type RequestPayload struct {
	Url string `json:"url"`
}

type ResponsePayload struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

type MappingURL struct {
	Url         string `json:"url"`
	GeneratedId string `json:"generatedId"`
}

func Connection(conn *sql.DB) {
	db = conn
}

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func Insert(mapping MappingURL) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into url_mapping (url, generated_id, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := db.QueryRowContext(ctx, stmt,
		mapping.Url,
		mapping.GeneratedId,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	log.Println(err)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
