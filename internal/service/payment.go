package service

import (
	"context"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/storage"
)

type PaymentService struct {
	repo   *storage.Manager
	logger logger.RequestLogger
}

func NewPaymentService(repo *storage.Manager, logger logger.RequestLogger) *PaymentService {
	return &PaymentService{repo: repo, logger: logger}
}

func (s *PaymentService) Create(ctx context.Context, payment model.Payment) (uint, error) {
	payment.Status = model.StatusCreated

	return s.repo.Payment.Create(ctx, payment)
}

func (s *PaymentService) GetByID(ctx context.Context, ID uint) (*model.Payment, error) {
	return s.repo.Payment.GetByID(ctx, ID)
}

func (s *PaymentService) Update(ctx context.Context, payment model.Payment) error {
	return s.repo.Payment.Update(ctx, payment)
}

func (s *PaymentService) GetPayments(ctx context.Context) (map[string][]model.Payment, error) {
	incomes, err := s.repo.Payment.GetByIncome(ctx)
	if err != nil {
		return nil, err
	}

	outcomes, err := s.repo.Payment.GetByOutcome(ctx)
	if err != nil {
		return nil, err
	}

	nonCompleted, err := s.repo.Payment.GetByNonCompleted(ctx)
	if err != nil {
		return nil, err
	}

	return map[string][]model.Payment{
		"nonCompleted": nonCompleted,
		"incomes":      incomes,
		"outcomes":     outcomes,
	}, nil
}
