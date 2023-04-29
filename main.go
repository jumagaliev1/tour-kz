package main

import (
	"context"
	"github.com/labstack/gommon/log"
	_ "tour-kz/docs"
	"tour-kz/internal/app"
	"tour-kz/internal/config"
	"tour-kz/internal/logger"
)

// @title Tour-KZ API
// @version 1.0

// @contact.name Alibi Zhumagaliyev
// @contact.url @AZhumagaliyev
// @contact.email alibi.zhumagaliyev@gmail.com

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description OAuth protects our entity endpoints
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cfg, err := config.New("configs/")
	if err != nil {
		log.Error(err)
	}
	a := app.New(cfg, logger.RequestLogger{})
	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
