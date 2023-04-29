package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID  uint
	Balance int
	User    User `gorm:"foreignKey:UserID"`
}

type AddBalanceReq struct {
	UserID uint `json:"user_id"`
	Amount int  `json:"amount"`
}
