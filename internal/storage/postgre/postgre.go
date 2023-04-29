package postgre

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tour-kz/internal/config"
	"tour-kz/internal/model"
)

func Dial(ctx context.Context, cfg config.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.URI()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.User{}, &model.Referral{}, &model.Account{}, &model.Payment{})
	return db, nil
}
