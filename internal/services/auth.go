package services

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"borda/pkg/hash"

	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repo   repository.UserRepository
	hasher hash.PasswordHasher
}

func NewAuthService(repo repository.UserRepository, hasher hash.PasswordHasher) *AuthService {
	return &AuthService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *AuthService) SignUp(input domain.UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	_, err = s.repo.Create(input.Username, passwordHash, input.Contact)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			// 400
			return err
		}
		// 500
		return err
	}

	return nil
}

func (s *AuthService) SignIn(username, password string) (string, error) {

	user, err := s.repo.FindUser(username, password)
	if err != nil {
		return "", err
	}

	jwtConf := config.JWT()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jwtConf.ExpireTime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.Id),
	})

	// TODO: save token somewhere
	return token.SignedString([]byte(jwtConf.SigningKey))
}
