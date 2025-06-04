package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
	CredPath    string `env:"RDS_VAULT_CRED_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	DetailsPath string `env:"RDS_VAULT_DETAIL_PATH" envDefault:"secret/data/chewedfeed/details"`

	ExpireTime time.Time
}

type Details struct {
	Host              string        `env:"RDS_HOSTNAME" envDefault:"postgres.chewedfeed"`
	Port              int           `env:"RDS_PORT" envDefault:"5432"`
	User              string        `env:"RDS_USERNAME"`
	Password          string        `env:"RDS_PASSWORD"`
	DBName            string        `env:"RDS_DB" envDefault:"postgres"`
	RawURL            string        `env:"RDS_URL"`
	ConnectionTimeout time.Duration `env:"RDS_CONNECTION_TIMEOUT" envDefault:"10s"`
	ExtraParams       string
}

type System struct {
	Context context.Context

	Details

	VaultDetails
	VaultHelper *vaultHelper.VaultHelper
}

func NewSystem() *System {
	return &System{
		Context: context.Background(),
	}
}

func (s *System) Setup(vd VaultDetails, vh vaultHelper.VaultHelper) {
	s.VaultDetails = vd
	s.VaultHelper = &vh
}

func (s *System) Build() (*Details, error) {
	gen, err := s.buildGeneric()
	if err != nil {
		return nil, err
	}

	if s.VaultHelper != nil {
		return s.buildVault()
	}

	return gen, nil
}

func (s *System) buildGeneric() (*Details, error) {
	rds := &Details{}
	if err := env.Parse(rds); err != nil {
		return rds, logs.Errorf("failed to parse database env: %v", err)
	}

	s.Details = *rds

	if rds.RawURL != "" {
		if err := s.ParseConnectionString(rds.RawURL); err != nil {
			return nil, err
		}
	}

	return rds, nil
}

func vaultSecretError(key string, err error) error {
	if err != nil {
		if err.Error() != fmt.Sprintf("key: '%s' not found", key) {
			return nil
		}
		return logs.Errorf("failed to get %s: %v", key, err)
	}

	return nil
}

func (s *System) buildVault() (*Details, error) {
	rds := &Details{}
	vh := *s.VaultHelper

	// Get Credentials
	if err := vh.GetSecrets(s.VaultDetails.CredPath); err != nil {
		return rds, logs.Errorf("failed to get cred secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return rds, logs.Error("no rds cred secrets found")
	}

	if s.Details.User == "" {
		secret, err := vh.GetSecret("username")
		if err := vaultSecretError("username", err); err != nil {
			return nil, logs.Errorf("failed to get username: %v", err)
		}
		rds.User = secret
	} else {
		rds.User = s.Details.User
	}

	if s.Details.Password == "" {
		secret, err := vh.GetSecret("password")
		if err := vaultSecretError("password", err); err != nil {
			return nil, logs.Errorf("failed to get password: %v", err)
		}
		rds.Password = secret
	} else {
		rds.Password = s.Details.Password
	}

	// Get Details
	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return rds, logs.Errorf("failed to get local secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return rds, logs.Error("no rds detail secrets found")
	}

	// get the port based on the username, since port has a default in env
	if s.Details.User == "" && s.Details.Port == 5432 {
		secret, err := vh.GetSecret("rds-port")
		if err != nil {
			if err.Error() != fmt.Sprint("key: 'rds-port' not found") {
				return nil, logs.Errorf("failed to get port: %v", err)
			}
			secret = "5432"
		}
		if secret != "" {
			iport, err := strconv.Atoi(secret)
			if err != nil {
				return nil, logs.Errorf("failed to parse port: %v", err)
			}
			rds.Port = iport
		}
	} else {
		rds.Port = s.Details.Port
	}

	// get the db based on the username, since db has a default in env
	if s.Details.User == "" {
		secret, err := vh.GetSecret("rds-db")
		if err := vaultSecretError("rds-db", err); err != nil {
			secret = "postgres"
		}
		rds.DBName = secret
	} else {
		rds.DBName = s.Details.DBName
	}

	// get the host based on the username, since host has a default in env
	if s.Details.User == "" {
		secret, err := vh.GetSecret("rds-hostname")
		if err := vaultSecretError("rds-hostname", err); err != nil {
			secret = "db.chewed-k8s.net"
		}
		rds.Host = secret
	} else {
		rds.Host = s.Details.Host
	}

	s.ExpireTime = time.Now().Add(time.Duration(vh.LeaseDuration()) * time.Second)
	s.Details = *rds

	return rds, nil
}

func (s *System) GetPGXClient(ctx context.Context) (*pgx.Conn, error) {
	if s.VaultHelper != nil && time.Now().Unix() > (s.VaultDetails.ExpireTime.Unix()-3600) {
		logs.Infof("vault expired, rebuilding, new expire time is %v", s.VaultDetails.ExpireTime)

		if _, err := s.buildVault(); err != nil {
			return nil, logs.Errorf("failed to build vault: %v", err)
		}
	}

	timeoutContext, cancel := context.WithTimeout(ctx, s.Details.ConnectionTimeout)
	s.Context = timeoutContext
	defer cancel()

	client, err := pgx.Connect(timeoutContext, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName, s.Details.ExtraParams))
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("failed to get db client: %v", err)
	}

	return client, nil
}

func (s *System) GetPGXPoolClient(ctx context.Context) (*pgxpool.Pool, error) {
	if s.VaultHelper != nil && time.Now().Unix() > (s.VaultDetails.ExpireTime.Unix()-3600) {
		logs.Infof("vault expired, rebuilding, new expire time is %v", s.VaultDetails.ExpireTime)

		if _, err := s.buildVault(); err != nil {
			return nil, logs.Errorf("failed to build vault: %v", err)
		}
	}

	timeoutContext, cancel := context.WithTimeout(ctx, s.Details.ConnectionTimeout)
	s.Context = timeoutContext
	defer cancel()

	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName, s.Details.ExtraParams))
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("failed to get db client: %v", err)
	}
	config.MaxConns = 10
	config.MaxConnIdleTime = s.Details.ConnectionTimeout

	client, err := pgxpool.NewWithConfig(timeoutContext, config)
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("failed to get db client: %v", err)
	}

	return client, nil
}

func (s *System) ClosePGX(ctx context.Context, conn *pgx.Conn) error {
	if err := conn.Close(ctx); err != nil {
		return logs.Errorf("failed to close db client: %v", err)
	}
	return nil
}

func (s *System) ParseConnectionString(connStr string) error {
	str, err := url.Parse(connStr)
	if err != nil {
		return logs.Errorf("failed to parse connection string: %v", err)
	}
	if str.Scheme != "postgres" && str.Scheme != "postgresql" {
		return logs.Errorf("invalid connection string scheme: %s", str.Scheme)
	}

	s.Details.Host = str.Hostname()
	if port, err := strconv.Atoi(str.Port()); err == nil {
		s.Details.Port = port
	} else {
		return logs.Errorf("invalid connection string port: %s", str.Port())
	}

	s.Details.User = str.User.Username()
	if pword, ok := str.User.Password(); ok {
		s.Details.Password = pword
	} else {
		return logs.Error("no password found in connection string")
	}

	if parts := strings.Split(str.Path, "/"); len(parts) > 1 {
		s.Details.DBName = parts[1]
	}

	if len(str.Query()) > 0 {
		s.Details.ExtraParams = str.RawQuery
	}

	return nil
}
