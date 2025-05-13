package payload

type AuthRegisterPayload struct {
	Username string  `json:"username" validate:"required,min=1,max=50"`
	Password string  `json:"password" validate:"required,min=1,max=12"`
	Role     *string `json:"role" validate:"max=20,oneof=user admin"`
}

type AuthLoginPayload struct {
	Username string `json:"username" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required,min=1,max=12"`
}

type AuthRegisterResponse struct {
	UserID  uint   `json:"userID"`
	Message string `json:"message"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}
