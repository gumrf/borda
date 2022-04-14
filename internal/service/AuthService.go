package service

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"borda/pkg/hash"
	"fmt"
	"strconv"
	"time"

	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	hasher   hash.PasswordHasher
}

func NewAuthService(ur repository.UserRepository, tr repository.TeamRepository,
	hasher hash.PasswordHasher) *AuthService {

	return &AuthService{
		userRepo: ur,
		teamRepo: tr,
		hasher:   hasher,
	}
}

func (s *AuthService) SignUp(input domain.SignUpInput) error {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	if _, err := s.userRepo.SaveUser(input.Username, hashedPassword, input.Contact); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}
		return err
	}

	return nil
}

func (s *AuthService) SignIn(input domain.SignInInput) (string, error) {
	// Hash password
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	// Find user with username and password
	user, err := s.userRepo.GetUserByCredentials(input.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	// Get user role
	role, err := s.userRepo.GetUserRole(user.Id)
	if err != nil {
		return "", err
	}

	jwtConf := config.JWT()

	fmt.Println(jwtConf)

	// Generate JWT token
	claims := jwt.MapClaims{
		"iss":   strconv.Itoa(user.Id),
		"exp":   jwt.NewNumericDate(time.Now().Add(jwtConf.ExpireTime)),
		"iat":   jwt.NewNumericDate(time.Now()),
		"aud":   "borda-v1",
		"scope": []string{role.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtConf.SigningKey))
}

func (s *AuthService) VerifyUserTeam(userId int) (int, bool) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return 0, false
	}

	if user.TeamId <= 0 {
		return 0, false
	}

	return user.TeamId, true
}
