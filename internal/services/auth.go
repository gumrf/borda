package services

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO: Service for generating statistics based on competition results (image or pdf)

type AuthService struct {
	repo *repository.Repository
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(r *repository.Repository) *AuthService {
	return &AuthService{
		repo: r,
	}
}

func (s *AuthService) SignUp(user domain.User) (int, error) {

}

func (s *AuthService) SignIn(username, password string) (string, error) {
	key, hours := config.GetJwtEntity()

	user, err := s.repo.Users.FindUser(username, password)
	if err != nil {
		return "", err
	}

	id := user.Id

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: id,
	})

	return token.SignedString([]byte(key))

}
