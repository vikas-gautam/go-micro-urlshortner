package models

import (
	"time"
)

//#################################################################################

//creating models against database tables

type Users struct {
	ID         int
	First_name string    `json:"first_name" validate:"required,min=2,max=30"`
	Last_name  string    `json:"last_name" validate:"required,min=2,max=30"`
	Email      string    `json:"email"      validate:"email,required"`
	Password   string    `json:"password"   validate:"required,min=6"`
	Status     string    `json:"status"`
	Phone      string    `json:"phone"      validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updtaed_at"`
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

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
