package data

import (
	"context"
	"database/sql"
	"fmt"
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
	Source_ip    string
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

	log.Println(mapping.Source_ip)

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

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func InsertURLClickCounter(mapping URLClickCounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int

	query := `insert into url_click_counter (short_url, click_counter, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := db.QueryRowContext(ctx, query,
		mapping.ShortURL,
		mapping.ClickCounter,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// update counter value as per query
func UpdateURLClickCounter(short_url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var mapping URLClickCounter
	mapping.ShortURL = short_url

	query := `select click_counter from url_click_counter where short_url=$1`

	err := db.QueryRowContext(ctx, query,
		mapping.ShortURL,
	).Scan(&mapping.ClickCounter)

	if err == sql.ErrNoRows {
		log.Println("No rows has been matched for given shortURL", err)
		return fmt.Errorf("no matching row found for short_url=%s", mapping.ShortURL)
	} else if err != nil {
		return err
	}

	// Update the ClickCounter column
	_, err = db.ExecContext(ctx, "UPDATE url_click_counter SET click_counter = click_counter + 1 WHERE short_url=$1", mapping.ShortURL)
	if err != nil {
		return err
	}

	return nil
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
