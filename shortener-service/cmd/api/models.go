package main

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
