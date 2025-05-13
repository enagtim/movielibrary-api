package service

import (
	"auth-service/internal/payload"
	"auth-service/internal/repository"
	"auth-service/pkg/consts"
	"auth-service/pkg/jwt"
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

func (s *AuthService) Login(ctx context.Context, p *payload.AuthLoginPayload) (string, error) {
	user, err := s.AuthRepository.GetUserByUsername(ctx, p.Username)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p.Password))

	if err != nil {
		return "", consts.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, user.Role)

	if err != nil {
		return "", consts.ErrGenerateToken
	}

	return token, nil

}
