package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func Connection(conn *sql.DB) {
	db = conn
}

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
	Url    string `json:"url"`
	Domain string `json:"domain"`
}

type ResponsePayload struct {
	ActualURL string `json:"actualurl"`
	ShortUrl  string `json:"shortUrl"`
}

type MappingURL struct {
	Url         string `json:"url"`
	GeneratedId string `json:"generatedId"`
}

type RowFromDB struct {
	ID           int
	Url          string
	generated_id string
	created_at   time.Time
	updated_at   time.Time
}

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func InsertUrl(mapping MappingURL) (string, error) {
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
		return "", err
	}

	return mapping.GeneratedId, nil
}

func GetUrlByid(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, url, generated_id, created_at, updated_at from url_mapping where generated_id = $1`

	var data RowFromDB

	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&data.ID,
		&data.Url,
		&data.generated_id,
		&data.created_at,
		&data.updated_at,
	)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return data.Url, nil
}
