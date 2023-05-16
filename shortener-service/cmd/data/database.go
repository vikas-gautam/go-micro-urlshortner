package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/vikas-gautam/go-micro-urlshortner/shortener-service/cmd/models"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func Connection(conn *sql.DB) {
	db = conn
}

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func InsertUrl(mapping models.URLMapping) (string, error) {
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
func InsertURLClickCounter(mapping models.URLClickCounter) error {
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

	var mapping models.URLClickCounter
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

	var data models.URLMapping

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

func GetCounterByshortUrl(short_url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, short_url, click_counter, created_at, updated_at from url_click_counter where short_url = $1`

	var data models.URLClickCounter

	row := db.QueryRowContext(ctx, query, short_url)

	err := row.Scan(
		&data.ID,
		&data.ShortURL,
		&data.ClickCounter,
		&data.Created_at,
		&data.Updated_at,
	)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return data.ClickCounter, nil
}

func CheckSourceIpExistence(source_ip string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select counter from url_generate_restrictions where source_ip=$1`

	var data models.URLGenerateRestrictions

	row := db.QueryRowContext(ctx, query, source_ip)

	err := row.Scan(
		&data.Counter,
	)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return data.Counter, nil

}

// Insert inserts a new mappingURL into the database, and returns the ID of the newly inserted row
func InsertURLGenerateRestrictions(mapping models.URLGenerateRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int

	query := `insert into url_generate_restrictions (source_ip, status, counter, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id`

	err := db.QueryRowContext(ctx, query,
		mapping.Source_ip,
		mapping.Status,
		mapping.Counter,
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
func UpdateURLGenerateRestrictions(source_ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var mapping models.URLGenerateRestrictions
	mapping.Source_ip = source_ip

	query := `select counter from url_generate_restrictions where source_ip=$1`

	err := db.QueryRowContext(ctx, query,
		mapping.Source_ip,
	).Scan(&mapping.Counter)

	if err == sql.ErrNoRows {
		log.Println("No rows has been matched for given source_ip", err)
		var mappingRestrictions models.URLGenerateRestrictions
		mappingRestrictions.Source_ip = source_ip
		mappingRestrictions.Status = "Active"
		mappingRestrictions.Counter = 1

		err = InsertURLGenerateRestrictions(mappingRestrictions)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// Update the ClickCounter column
	_, err = db.ExecContext(ctx, "UPDATE url_generate_restrictions SET counter = counter + 1 WHERE source_ip=$1", mapping.Source_ip)
	if err != nil {
		return err
	}

	return nil
}

// update counter value as per query
func UpdateStatusURLGenerateRestrictions(source_ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var mapping models.URLGenerateRestrictions
	mapping.Source_ip = source_ip

	// Update the status column
	_, err := db.ExecContext(ctx, "UPDATE url_generate_restrictions SET status= 'Blocked' WHERE source_ip=$1", mapping.Source_ip)
	if err != nil {
		return err
	}

	return nil
}
