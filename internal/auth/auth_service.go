package auth

import (
	"borda/internal/config"
	"borda/internal/pkg/core"
	"borda/pkg/hash"
	"fmt"

	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepository core.UserRepository
	teamRepository core.TeamRepository
	hasher         hash.PasswordHasher
	jwtConfig      config.JWTConfig
}

func NewAuthService(uRepo core.UserRepository, tRepo core.TeamRepository,
	hasher hash.PasswordHasher, jwtConfig config.JWTConfig) *AuthService {

	return &AuthService{
		userRepository: uRepo,
		teamRepository: tRepo,
		hasher:         hasher,
		jwtConfig:      jwtConfig,
	}
}

func (as *AuthService) Register(request RegistrationRequest) error {
	hashedPassword, err := as.hasher.Hash(request.Password)
	if err != nil {
		return err
	}

	user := core.User{
		Username: request.Username,
		Password: hashedPassword,
		Contact:  request.Password,
	}

	if _, err := as.userRepository.Save(context.Background(), user); err != nil {
		// if errors.Is(err, domain.ErrUserAlreadyExists) {
		// 	return err
		// }
		return err
	}

	return nil
}

type AuthorizationResponse struct {
	Token string `json:"token"`
}

func (as *AuthService) Authorize(credentials AuthorizationRequest) (AuthorizationResponse, error) {
	hashedPassword, err := as.hasher.Hash(credentials.Password)
	if err != nil {
		return AuthorizationResponse{}, err
	}

	// Find user with username and password
	user, err := as.userRepository.FindByCredentials(credentials.Username, hashedPassword)
	if err != nil {
		return AuthorizationResponse{}, err
	}

	// TODO
	// Get user role
	// role, err := s.userRepo.GetUserRole(user.Id)
	// if err != nil {
	// 	return "", err
	// }

	role := core.Role{
		Id:   1,
		Name: "admin",
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"iss":   strconv.Itoa(user.Id),
		"exp":   jwt.NewNumericDate(time.Now().Add(as.jwtConfig.ExpireTime)),
		"iat":   jwt.NewNumericDate(time.Now()),
		"aud":   "borda-v1",
		"scope": []string{role.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(as.jwtConfig.SigningKey)

	signedToken, err := token.SignedString([]byte(as.jwtConfig.SigningKey))
	if err != nil {
		return AuthorizationResponse{}, nil
	}

	return AuthorizationResponse{Token: signedToken}, nil
}

// func (s *AuthService) VerifyUserTeam(userId int) (int, bool) {
// 	user, err := s.userRepo.FindById(userId)
// 	if err != nil {
// 		return 0, false
// 	}

// 	if user.TeamId <= 0 {
// 		return 0, false
// 	}

// 	return user.TeamId, true
// }
