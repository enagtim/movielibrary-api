package payload

import "time"

type ActorPaylod struct {
	Name      string    `json:"name" validate:"required,min=1,max=150"`
	Gender    string    `json:"gender" validate:"required, oneof=male female other"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
}

type PartialUpdateActorPaylod struct {
	Name      *string    `json:"name" validate:"min=1,max=150"`
	Gender    *string    `json:"gender" validate:"oneof=male female other"`
	BirthDate *time.Time `json:"birth_date"`
}
