package usecase

type AdminUsecase struct{}

// Verify interface compliance
var _ AdminUsecaseI = (*AdminUsecase)(nil)