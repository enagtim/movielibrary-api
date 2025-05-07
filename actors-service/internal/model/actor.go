package model

import "time"

type Actor struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}
