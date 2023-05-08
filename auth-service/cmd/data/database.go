package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/vikas-gautam/go-micro-urlshortner/auth-service/cmd/models"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func Connection(conn *sql.DB) {
	db = conn
}

// Insert inserts a new users's signup data into the database, and returns the ID of the newly inserted row
func InsertSignupData(data models.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (first_name, last_name, email, password, status, phone, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	err := db.QueryRowContext(ctx, stmt,
		data.First_name,
		data.Last_name,
		data.Email,
		data.Password,
		data.Status,
		data.Phone,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	log.Println(err)

	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmailid(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, first_name, last_name, email, password, status, phone, created_at, updated_at from users where email = $1`

	var data models.Users

	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&data.ID,
		&data.First_name,
		&data.Last_name,
		&data.Email,
		&data.Password,
		&data.Status,
		&data.Phone,
		&data.Created_at,
		&data.Updated_at,
	)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func GetUserByPhone(phone string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, first_name, last_name, email, password, status, phone, created_at, updated_at from users where phone = $1`

	var data models.Users

	row := db.QueryRowContext(ctx, query, phone)

	err := row.Scan(
		&data.ID,
		&data.First_name,
		&data.Last_name,
		&data.Email,
		&data.Password,
		&data.Status,
		&data.Phone,
		&data.Created_at,
		&data.Updated_at,
	)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func GetHashedPasswordByEmailid(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, first_name, last_name, email, password, status, phone, created_at, updated_at from users where email = $1`

	var data models.Users

	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&data.ID,
		&data.First_name,
		&data.Last_name,
		&data.Email,
		&data.Password,
		&data.Status,
		&data.Phone,
		&data.Created_at,
		&data.Updated_at,
	)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return data.Password, nil
}
