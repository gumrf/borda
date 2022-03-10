package services

import "borda/internal/repository"

// TODO: Service for generating statistics based on competition results (image or pdf)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(r *repository.Repository) *AuthService {
	return &AuthService{
		repo: r,
	}
}
