package postgres

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Details struct {
	Host              string        `env:"RDS_HOSTNAME" envDefault:"postgres.chewedfeed" vault:"rds-hostname"`
	Port              int           `env:"RDS_PORT" envDefault:"5432" vault:"rds-port"`
	User              string        `env:"RDS_USERNAME" vault:"username,creds"`
	Password          string        `env:"RDS_PASSWORD" vault:"password,creds"`
	DBName            string        `env:"RDS_DB" envDefault:"postgres" vault:"rds-db"`
	RawURL            string        `env:"RDS_URL"`
	ConnectionTimeout time.Duration `env:"RDS_CONNECTION_TIMEOUT" envDefault:"10s"`
	ExtraParams       string
}

type System struct {
	Context context.Context

	Details
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Build() (*Details, error) {
	return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
	rds := &Details{}
	if err := env.Parse(rds); err != nil {
		return rds, logs.Errorf("postgres: unable to parse env: %v", err)
	}

	s.Details = *rds

	if rds.RawURL != "" {
		if err := s.ParseConnectionString(rds.RawURL); err != nil {
			return nil, err
		}
	}

	return rds, nil
}

func (s *System) GetPGXClient(ctx context.Context) (*pgx.Conn, error) {
	timeoutContext, cancel := context.WithTimeout(ctx, s.Details.ConnectionTimeout)
	defer cancel()

	client, err := pgx.Connect(timeoutContext, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName, s.Details.ExtraParams))
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("postgres: unable to connect: %v", err)
	}

	return client, nil
}

// GetPGXPoolClient opens a pool to the database. Each opt can change the pool's config after go-config's
// defaults (MaxConns 10, idle timeout = ConnectionTimeout) and before the pool opens: pgx's hooks, such as
// PrepareConn to set a session setting from the request's context, or different limits.
func (s *System) GetPGXPoolClient(ctx context.Context, opts ...func(*pgxpool.Config)) (*pgxpool.Pool, error) {
	timeoutContext, cancel := context.WithTimeout(ctx, s.Details.ConnectionTimeout)
	defer cancel()

	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName, s.Details.ExtraParams))
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("postgres: unable to parse pool config: %v", err)
	}
	config.MaxConns = 10
	config.MaxConnIdleTime = s.Details.ConnectionTimeout
	for _, opt := range opts {
		opt(config)
	}

	client, err := pgxpool.NewWithConfig(timeoutContext, config)
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("postgres: unable to create pool: %v", err)
	}

	return client, nil
}

func (s *System) ClosePGX(ctx context.Context, conn *pgx.Conn) error {
	if err := conn.Close(ctx); err != nil {
		return logs.Errorf("postgres: unable to close connection: %v", err)
	}
	return nil
}

func (s *System) ParseConnectionString(connStr string) error {
	str, err := url.Parse(connStr)
	if err != nil {
		return logs.Errorf("postgres: unable to parse connection string: %v", err)
	}
	if str.Scheme != "postgres" && str.Scheme != "postgresql" {
		return logs.Errorf("postgres: unable to use connection string scheme: %s", str.Scheme)
	}

	s.Details.Host = str.Hostname()
	if port, err := strconv.Atoi(str.Port()); err == nil {
		s.Details.Port = port
	} else {
		return logs.Errorf("postgres: unable to parse connection string port: %s", str.Port())
	}

	s.Details.User = str.User.Username()
	if pword, ok := str.User.Password(); ok {
		s.Details.Password = pword
	} else {
		return logs.Error("postgres: unable to find password in connection string")
	}

	if parts := strings.Split(str.Path, "/"); len(parts) > 1 {
		s.Details.DBName = parts[1]
	}

	if len(str.Query()) > 0 {
		s.Details.ExtraParams = str.RawQuery
	}

	return nil
}
