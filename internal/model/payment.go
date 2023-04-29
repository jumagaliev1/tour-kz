package model

import "gorm.io/gorm"

const (
	StatusCreated  = "CREATED"
	StatusComplete = "COMPLETED"
	StatusCancel   = "CANCELED"

	TypeIncome  = "INCOME"
	TypeOutcome = "OUTCOME"
)

type Payment struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
	Type   string `json:"type"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TypeCreateReq struct {
	UserID uint   `json:"user_id"`
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

func (p *TypeCreateReq) ToPayment() Payment {
	return Payment{
		UserID: p.UserID,
		Amount: p.Amount,
		Type:   p.Type,
	}
}
