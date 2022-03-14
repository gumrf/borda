package services

import (
	"borda/internal/config"
	"borda/internal/domain"
	"borda/internal/repository"
	"borda/pkg/hash"
	"fmt"

	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func (s *AuthService) SignUp(input domain.UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	// TODO:
	//		Attach user to the team.
	//		If parsing token or creating new team fails, user should't be created.
	// 		To achive prosess  should be run in transaction.
	userId, err := s.userRepo.CreateNewUser(input.Username, passwordHash, input.Contact)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}
		return err
	}

	switch input.AttachTeamMethod {
	case "create":
		_, err = s.teamRepo.CreateNewTeam(userId, input.AttachTeamAttribute)
		if err != nil {
			return err
		}
	case "join":
		team, err := s.teamRepo.GetTeamByToken(input.AttachTeamAttribute)
		if err != nil {
			return err
		}

		err = s.teamRepo.AddMember(team.Id, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AuthService) SignIn(input domain.UserSignInInput) (string, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	fmt.Println(passwordHash)

	user, err := s.userRepo.FindUserByCredentials(input.Username, passwordHash)
	if err != nil {
		return "", err
	}

	jwtConf := config.JWT()

	fmt.Println(jwtConf)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jwtConf.ExpireTime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.Id),
	})

	// TODO: save token somewhere
	return token.SignedString([]byte(jwtConf.SigningKey))
}
