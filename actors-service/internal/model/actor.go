package model

import "time"

type Actor struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type ActorWithMovies struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	Movies    []string  `json:"movies"`
}
