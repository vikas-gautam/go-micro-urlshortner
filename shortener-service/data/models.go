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

//#################################################################################

//creating models against database tables

type Users struct {
	ID         int
	First_name string
	Last_name  string
	Email      string
	Password   string
	Status     string
	Created_at time.Time
	Updated_at time.Time
}

type URLMapping struct {
	ID           int
	Source_ip    int
	Generated_id string
	Url          string
	User_type    string
	Email        string
	Created_at   time.Time
	Updated_at   time.Time
}

type UserType struct {
	ID         int
	Type       int
	Created_at time.Time
	Updated_at time.Time
}

type URLGenerateRestrictions struct {
	ID        int
	Source_ip int
	Status    string
	Counter   int
}

type URLClickCounter struct {
	ID           int
	ShortURL     string
	ClickCounter int
}

//#########################################################################

type HealthPayload struct {
	Message string
	Status  int
}

// To take the user's Post request #############################################
type RequestPayload struct {
	Url string `json:"url"`
}

type ResponsePayload struct {
	ActualURL string `json:"actualurl"`
	ShortUrl  string `json:"shortUrl"`
}

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func InsertUrl(mapping URLMapping) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into url_mapping (source_ip, generated_id, url, user_type, email, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := db.QueryRowContext(ctx, stmt,
		mapping.Source_ip,
		mapping.Generated_id,
		mapping.Url,
		mapping.User_type,
		mapping.Email,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	log.Println(err)

	if err != nil {
		return "", err
	}

	return mapping.Generated_id, nil
}

func GetUrlByid(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, source_ip, generated_id, url, user_type, email, created_at, updated_at from url_mapping where generated_id = $1`

	var data URLMapping

	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&data.ID,
		&data.Source_ip,
		&data.Generated_id,
		&data.Url,
		&data.User_type,
		&data.Email,
		&data.Created_at,
		&data.Updated_at,
	)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return data.Url, nil
}
