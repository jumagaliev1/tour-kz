package app

import (
	"context"
	"tour-kz/internal/config"
	"tour-kz/internal/logger"
	"tour-kz/internal/service"
	"tour-kz/internal/storage"
	http "tour-kz/internal/transport"
	"tour-kz/internal/transport/http/handler"
	"tour-kz/internal/transport/middleware"
)

type App struct {
	cfg    *config.Config
	logger logger.RequestLogger
}

func New(cfg *config.Config, logger logger.RequestLogger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	stg, err := storage.New(ctx, a.cfg, a.logger)
	if err != nil {
		return err
	}

	svc, err := service.New(stg, a.cfg, a.logger)
	if err != nil {
		return err
	}

	jwtAuth := middleware.NewJWTAuth(a.cfg, svc.User)
	permission := middleware.NewPermission(svc.User)

	ctrl, err := handler.New(svc, jwtAuth, a.logger)
	if err != nil {
		return err
	}

	HTTPServer := http.NewServer(a.cfg, ctrl, jwtAuth, permission)

	return HTTPServer.StartHTTPServer(ctx)
}
