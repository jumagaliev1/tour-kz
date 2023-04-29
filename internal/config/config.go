package config

import (
	"fmt"
	"github.com/subosito/gotenv"
	"net"
	"os"
	"time"
)

const (
	defaultHTTPPort     = 8080
	defaultShutDownTime = 10 * time.Second
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	GRPC            string
}
type PostgresConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	SSLMode     string
	PingTimeout time.Duration
}
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	JWTKey   string
}

func (c PostgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.Username,
		c.Password,
		net.JoinHostPort(c.Host, c.Port),
		c.Database,
		c.SSLMode,
	)
}

func Init(cfg *Config) {
	cfg.Server.Host = os.Getenv("SERVER_HOST")

	cfg.Postgres.Username = os.Getenv("POSTGRES_USERNAME")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Database = os.Getenv("POSTGRES_DATABASE")
	cfg.Postgres.SSLMode = os.Getenv("POSTGRES_SSL_MODE")
}

func defaultInit(cfg *Config) {
	cfg.Server.Port = defaultHTTPPort
	cfg.Server.ShutdownTimeout = defaultShutDownTime
}

func New(path string) (*Config, error) {
	cfg := &Config{}

	defaultInit(cfg)

	if err := gotenv.Load(path + "/.env"); err != nil {
		return nil, err
	}

	Init(cfg)

	return cfg, nil
}
