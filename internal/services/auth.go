package services

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"borda/pkg/hash"
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
	var id int
	salt := config.GetPasswordSalt()

	MyHasher := hash.NewSHA1Hasher(salt)
	err := s.repo.Users.FindUserByUsename(user.Username)
	if err != nil {
		pswd, _ := MyHasher.Hash(user.Password)
		id, err = s.repo.Users.Create(user.Username, pswd, user.Contact)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
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
