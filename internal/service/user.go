package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
	"tour-kz/internal/config"
	"tour-kz/internal/logger"
	"tour-kz/internal/model"
	"tour-kz/internal/storage"
)

type UserService struct {
	repo   *storage.Manager
	cfg    *config.Config
	logger logger.RequestLogger
}

func NewUserService(repo *storage.Manager, cfg *config.Config, logger logger.RequestLogger) *UserService {
	return &UserService{repo: repo, cfg: cfg, logger: logger}
}

func (s *UserService) Create(ctx context.Context, input model.UserCreateReq) (*model.User, error) {
	var err error

	input.Password, err = s.HashPassword(input.Password)
	if err != nil {
		s.logger.Logger(ctx).Error(err)
		return nil, err
	}

	var parent *model.User
	if input.ReferralCode != "" {
		parent, err = s.repo.User.GetByReferralCode(ctx, input.ReferralCode)
		if err != nil {
			return nil, err
		}
	}

	input.ReferralCode = generateReferralCode()

	u, err := s.repo.User.Create(ctx, *input.MapperToUser())
	if err != nil {
		return nil, err
	}
	if parent != nil {
		_, err = s.repo.Referral.Create(ctx, model.Referral{UserID: u.ID, ParentID: parent.ID})
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *UserService) Update(ctx context.Context, user model.User) error {
	return s.repo.User.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, ID int) error {
	return s.repo.User.Delete(ctx, ID)
}

func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.repo.User.GetAll(ctx)
}

func (s *UserService) CheckPassword(encPass, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encPass), []byte(providedPassword))
}

func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func (s *UserService) Auth(ctx context.Context, user model.AuthUser) error {
	userFromDB, userErr := s.repo.User.GetByEmail(ctx, user.Email)
	if userErr != nil {
		s.logger.Logger(ctx).Error(userErr)
		return userErr
	}
	s.logger.Logger(ctx).Info(s.HashPassword(user.Password))
	s.logger.Logger(ctx).Info(userFromDB.Password)
	checkErr := s.CheckPassword(userFromDB.Password, user.Password)
	if checkErr != nil {
		s.logger.Logger(ctx).Error(checkErr.Error())
		return checkErr
	}

	return nil
}

func (s *UserService) RefreshToken() (string, error) {
	b := make([]byte, 32)

	str := rand.NewSource(time.Now().Unix())
	r := rand.New(str)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *UserService) GenerateToken(user model.AuthUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.JWTClaim{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString([]byte(s.cfg.JWTKey))
}

func (s *UserService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.cfg.JWTKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*model.JWTClaim)
	if !ok {
		return "", errors.New("token claims are not of type *tokeClaims*")
	}

	return claims.Email, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.User.GetByEmail(ctx, email)
}

func (s *UserService) GetUserFromRequest(ctx context.Context) (*model.User, error) {
	email, ok := ctx.Value(model.ContextEmail).(string)
	if !ok {
		s.logger.Logger(ctx).Error("not valid context username")
		user := &model.User{
			FirstName: "Anonymous",
			LastName:  "Guest",
			Email:     "guest@gmail.com",
		}
		return user, nil
	}

	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//func (s *UserService) ChangePassword(ctx context.Context, body model.PasswordReq) error {
//	user, err := s.GetUserFromRequest(ctx)
//	if err != nil {
//		s.logger.Logger(ctx).Error(err)
//		return err
//	}
//	checkErr := s.CheckPassword(user.Password, body.OldPassword)
//	if checkErr != nil {
//		s.logger.Logger(ctx).Error(checkErr)
//		return checkErr
//	}
//
//	hash, err := s.HashPassword(body.Password)
//	if err != nil {
//		s.logger.Logger(ctx).Error(err)
//		return err
//	}
//
//	user.Password = hash
//
//	return s.Update(ctx, *user)
//}

func generateReferralCode() string {
	// Generate a random string
	bytes := make([]byte, 8)
	rand.Read(bytes)
	randomString := base64.URLEncoding.EncodeToString(bytes)

	// Get the current timestamp
	timestamp := time.Now().UnixNano()

	// Combine the random string and the timestamp to generate the referral code
	referralCode := fmt.Sprintf("%s:%d", randomString, timestamp)

	return referralCode
}
