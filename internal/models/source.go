package models

type Source struct {
	Email     Email  `json:"email"`
	Timestamp string `json:"@timestamp"`
	Path      string `json:"path"`
}
