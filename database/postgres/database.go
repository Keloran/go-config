package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	Host     string `env:"RDS_HOSTNAME" envDefault:"postgres.chewedfeed"`
	Port     int    `env:"RDS_PORT" envDefault:"5432"`
	User     string `env:"RDS_USERNAME"`
	Password string `env:"RDS_PASSWORD"`
	DBName   string `env:"RDS_DB" envDefault:"postgres"`
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

	return rds, nil
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
		if err != nil {
			return nil, logs.Errorf("failed to get username: %v", err)
		}
		rds.User = secret
	} else {
		rds.User = s.Details.User
	}

	if s.Details.Password == "" {
		secret, err := vh.GetSecret("password")
		if err != nil {
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
	if s.Details.User == "" {
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
		if err != nil {
			if err.Error() != fmt.Sprint("key: 'rds-db' not found") {
				return nil, logs.Errorf("failed to get db: %v", err)
			}
			secret = "postgres"
		}
		rds.DBName = secret
	} else {
		rds.DBName = s.Details.DBName
	}

	// get the host based on the username, since host has a default in env
	if s.Details.User == "" {
		secret, err := vh.GetSecret("rds-hostname")
		if err != nil {
			if err.Error() != fmt.Sprint("key: 'rds-hostname' not found") {
				return nil, logs.Errorf("failed to get host: %v", err)
			}
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

	client, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName))
	if err != nil {
		if strings.Contains(err.Error(), "operation was canceled") {
			return nil, err
		}
		return nil, logs.Errorf("failed to get db client: %v", err)
	}

	return client, nil
}

func (s *System) ClosePGX(ctx context.Context, conn pgx.Conn) error {
	if err := conn.Close(ctx); err != nil {
		return logs.Errorf("failed to close db client: %v", err)
	}

	return nil
}
