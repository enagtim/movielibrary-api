package model

import "time"

type Movie struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Decription  string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
}
