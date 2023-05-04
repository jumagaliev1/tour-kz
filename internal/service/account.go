package service

import (
	"context"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/storage"
)

type AccountService struct {
	repo   *storage.Manager
	logger logger.RequestLogger
}

func NewAccountService(repo *storage.Manager, logger logger.RequestLogger) *AccountService {
	return &AccountService{repo: repo, logger: logger}
}

func (s *AccountService) Create(ctx context.Context, account model.Account) (uint, error) {
	return s.repo.Account.Create(ctx, account)
}

func (s *AccountService) GetByUser(ctx context.Context, userID uint) (*model.Account, error) {
	return s.repo.Account.GetByUserID(ctx, userID)
}

func (s *AccountService) Update(ctx context.Context, account model.Account) error {
	return s.repo.Account.Update(ctx, account)
}

func (s *AccountService) UpdateLevels(ctx context.Context, userID uint, amount int) error {
	return s.repo.Account.UpdateBalance(ctx, userID, amount)
}
