package usecase

type TeamLeaderUsecase struct{}

// Verify interface compliance
var _ TeamLeaderUsecaseI = (*TeamLeaderUsecase)(nil)