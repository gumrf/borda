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

//func (s *AuthService) verifyData(input domain.SignUpInput) error {
//
//	// Реализована внутри функции SaveUser
//	// Проверка на имя пользователя
//	//err := s.userRepo.IsUsernameExists(input.Username)
//	//if err != nil {
//	//	return err
//	//}
//
//	//Проверка на имя команды или uid
//	switch input.AttachTeamMethod {
//	case "create":
//		//err = s.teamRepo.IsTeamNameExists(input.AttachTeamAttribute)
//		//if err != nil {
//		//	return err
//		//}
//	case "join":
//		// err := s.teamRepo.IsTeamTokenExists(input.AttachTeamAttribute)
//		// if err != nil {
//		// 	return err
//		// }
//
//		team, err := s.teamRepo.GetTeamByToken(input.AttachTeamAttribute)
//		if err != nil {
//			return err
//		}
//
//		err = s.teamRepo.IsTeamFull(team.Id)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (s *AuthService) SignUp(input domain.SignUpInput) error {
	//err := s.verifyData(input)
	//if err != nil {
	//	return err
	//}

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

	//switch input.AttachTeamMethod {
	//case "create":
	//	if _, err := s.teamRepo.SaveTeam(userId, input.AttachTeamAttribute); err != nil {
	//		return err
	//	}
	//case "join":
	//	team, err := s.teamRepo.GetTeamByToken(input.AttachTeamAttribute)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if err := s.teamRepo.AddMember(team.Id, userId); err != nil {
	//		return err
	//	}
	//}

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
