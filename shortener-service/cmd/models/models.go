package models

import (
	"database/sql"
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
	ID         int
	Source_ip  string
	Status     string
	Counter    int
	Created_at time.Time
	Updated_at time.Time
}

type URLClickCounter struct {
	ID           int
	ShortURL     string
	ClickCounter int
	Created_at   time.Time
	Updated_at   time.Time
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

type ResponsePayloadError struct {
	Status  int
	Message string
}

type ResponsePayloadUrlCounter struct {
	ShortUrl       string `json:"shortUrl"`
	TotalURLClicks int    `json:"totalURLClicks"`
}
