package transport

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
	"tour-kz/internal/config"
	"tour-kz/internal/transport/http/handler"
	"tour-kz/internal/transport/middleware"
)

type Server struct {
	cfg        *config.Config
	App        *echo.Echo
	handler    *handler.Manager
	jwt        *middleware.JWTAuth
	permission *middleware.Permission
	//middleware middleware.Middleware
}

func NewServer(cfg *config.Config, handler *handler.Manager, jwt *middleware.JWTAuth, permission *middleware.Permission) *Server {
	return &Server{
		cfg:        cfg,
		handler:    handler,
		jwt:        jwt,
		permission: permission,
	}
}

func (s *Server) StartHTTPServer(ctx context.Context) error {
	s.App = s.BuildEngine()
	s.SetupRoutes()

	go func() {
		if err := s.App.Start(fmt.Sprintf(":%v", s.cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:\v\n", err)
		}
	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.App.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:#{err}")
	}
	log.Print("server exited properly")
	return nil
}

func (s *Server) BuildEngine() *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	//e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			log.Info(map[string]interface{}{"URI": v.URI, "status": v.Status})
			return nil
		},
	}))

	//l := logger.Logger{context.Background()}

	return e
}
