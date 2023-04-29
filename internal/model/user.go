package model

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	Role         string `json:"role" gorm:"default:client"`
	Phone        string `json:"phone"`
	BankUUID     string `json:"bank_uuid"`
	ReferralCode string `json:"referral_code"`
}
type UserCreateReq struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	BankUUID     string `json:"bank_uuid"`
	ReferralCode string `json:"referral_code"`
}

func (u *UserCreateReq) MapperToUser() *User {
	return &User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		Password:     u.Password,
		ReferralCode: u.ReferralCode,
		Phone:        u.Phone,
		BankUUID:     u.BankUUID,
	}
}

var ContextEmail = contextKey("email")

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type contextKey string

type JWTClaim struct {
	Email          string
	StandardClaims jwt.StandardClaims
}

func (jwt *JWTClaim) Valid() error {
	return nil
}
