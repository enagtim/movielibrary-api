package service

import (
	"auth-service/internal/payload"
	"auth-service/internal/repository"
	"auth-service/pkg/consts"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, p *payload.AuthRegisterPayload) (uint, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)

	if err != nil {
		return 0, consts.ErrFailedHashedPassword
	}

	p.Password = string(hashedPassword)

	userID, err := s.AuthRepository.Register(ctx, p)

	if err != nil {
		return 0, err
	}

	return userID, nil

}
