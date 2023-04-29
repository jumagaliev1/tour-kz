package service

import (
	"context"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/storage"
)

type ReferralService struct {
	repo   *storage.Manager
	logger logger.RequestLogger
}

func NewReferralService(repo *storage.Manager, logger logger.RequestLogger) *ReferralService {
	return &ReferralService{repo: repo, logger: logger}
}

func (s *ReferralService) Create(ctx context.Context, referral model.Referral) (uint, error) {
	return s.repo.Referral.Create(ctx, referral)
}

func (s *ReferralService) GetFirstLine(ctx context.Context, userID uint) ([]model.User, error) {
	return s.repo.Referral.GetFirstLine(ctx, userID)
}

func (s *ReferralService) GetSecondLine(ctx context.Context, userID uint) ([]model.User, error) {
	return s.repo.Referral.GetSecondLine(ctx, userID)
}

func (s *ReferralService) GetThirdLine(ctx context.Context, userID uint) ([]model.User, error) {
	return s.repo.Referral.GetThirdLine(ctx, userID)
}

func (s *ReferralService) GetReferrals(ctx context.Context, userID uint) (map[string][]model.User, error) {
	line1, err := s.GetFirstLine(ctx, userID)
	if err != nil {
		return nil, err
	}
	line2, err := s.GetSecondLine(ctx, userID)
	if err != nil {
		return nil, err
	}
	line3, err := s.GetThirdLine(ctx, userID)
	if err != nil {
		return nil, err
	}

	return map[string][]model.User{"Line 1": line1, "Line 2": line2, "Line 3": line3}, nil

}
