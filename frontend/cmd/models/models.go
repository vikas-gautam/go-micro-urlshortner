package models

type URLCollection struct {
	ActualURL string
	ShortURL  string
}
type SuccessResponse struct {
	Response URLCollection
}

type AuthResponse struct {
	Status  int
	Message string
}

type Header struct {
	XForwardedFor string
	UserType      string
	Email         string
}
