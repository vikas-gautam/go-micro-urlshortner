package models

type URLCollection struct {
	ActualURL string
	ShortURL  string
}
type SuccessResponse struct {
	Response URLCollection
}
