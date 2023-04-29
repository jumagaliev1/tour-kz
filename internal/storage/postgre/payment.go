package postgre

import (
	"context"
	"gorm.io/gorm"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
)

type PaymentRepository struct {
	DB     *gorm.DB
	logger logger.RequestLogger
}

func NewPaymentRepository(DB *gorm.DB, logger logger.RequestLogger) *PaymentRepository {
	return &PaymentRepository{DB: DB, logger: logger}
}

func (r *PaymentRepository) Create(ctx context.Context, payment model.Payment) (uint, error) {
	if err := r.DB.WithContext(ctx).Create(&payment).Error; err != nil {
		return 0, err
	}

	return payment.ID, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, ID uint) (*model.Payment, error) {
	var payment *model.Payment
	if err := r.DB.WithContext(ctx).Find(&payment, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment model.Payment) error {
	if err := r.DB.WithContext(ctx).Save(&payment).Error; err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) GetByNonCompleted(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	rows, err := r.DB.WithContext(ctx).Raw(`
		SELECT p.ID, p.amount, p.status, p.type, p.created_at, 
		u.id, u.first_name, u.last_name, u.email, u.role, u.phone, u.referral_code 
		FROM payments p JOIN users u on u.id = p.user_id where p.status = ?`, model.StatusCreated).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var payment model.Payment
		rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Status,
			&payment.Type,
			&payment.CreatedAt,
			&payment.User.ID,
			&payment.User.FirstName,
			&payment.User.LastName,
			&payment.User.Email,
			&payment.User.Role,
			&payment.User.Phone,
			&payment.User.ReferralCode)
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetByIncome(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	rows, err := r.DB.WithContext(ctx).Raw(`
		SELECT p.ID, p.amount, p.status, p.type, p.created_at, 
		u.id, u.first_name, u.last_name, u.email, u.role, u.phone, u.referral_code 
		FROM payments p JOIN users u on u.id = p.user_id where p.status = ? and p.type = ?`, model.StatusCreated, model.TypeIncome).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var payment model.Payment
		rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Status,
			&payment.Type,
			&payment.CreatedAt,
			&payment.User.ID,
			&payment.User.FirstName,
			&payment.User.LastName,
			&payment.User.Email,
			&payment.User.Role,
			&payment.User.Phone,
			&payment.User.ReferralCode)
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetByOutcome(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	rows, err := r.DB.WithContext(ctx).Raw(`
		SELECT p.ID, p.amount, p.status, p.type, p.created_at, 
		u.id, u.first_name, u.last_name, u.email, u.role, u.phone, u.referral_code 
		FROM payments p JOIN users u on u.id = p.user_id where p.status = ? and p.type = ?`, model.StatusCreated, model.TypeOutcome).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var payment model.Payment
		rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Status,
			&payment.Type,
			&payment.CreatedAt,
			&payment.User.ID,
			&payment.User.FirstName,
			&payment.User.LastName,
			&payment.User.Email,
			&payment.User.Role,
			&payment.User.Phone,
			&payment.User.ReferralCode)
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) GetAll(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	if err := r.DB.WithContext(ctx).Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
