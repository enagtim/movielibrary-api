package payload

import (
	"actors-service/internal/model"
	"time"
)

type ActorPayload struct {
	Name      string    `json:"name" validate:"required,min=1,max=150"`
	Gender    string    `json:"gender" validate:"required, oneof=male female other"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
}

type PartialUpdateActorPayload struct {
	Name      *string    `json:"name" validate:"min=1,max=150"`
	Gender    *string    `json:"gender" validate:"oneof=male female other"`
	BirthDate *time.Time `json:"birth_date"`
}

type CreatedActorResponse struct {
	ActorID uint   `json:"actordID"`
	Message string `json:"message"`
}

type GetActorResponse struct {
	Data []model.ActorWithMovies `json:"data"`
}
