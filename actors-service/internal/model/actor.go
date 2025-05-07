package model

import "time"

type Actor struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" validate:"required,min=1,max=150"`
	Gender    string    `json:"gender" validate:"required, oneof=male female other"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
}
