package usecase

type UserUsecase struct{}

// Verify interface compliance
var _ UserUsecaseI = (*UserUsecase)(nil)