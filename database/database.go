package database

import (
  "context"
  "fmt"
  "strconv"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
	Address string
	Token   string

	CredPath   string `env:"RDS_VAULT_CRED_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	DetailPath string `env:"RDS_VAULT_DETAIL_PATH" envDefault:"secret/data/chewedfeed/details"`

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
  if s.VaultHelper != nil {
    return s.buildVault()
  }

  return s.buildGeneric()
}

func (s *System) buildGeneric() (*Details, error) {
  rds := &Details{}
  if err := env.Parse(rds); err != nil {
    return rds, logs.Errorf("failed to parse database env: %v", err)
  }

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
    return rds, logs.Error("no rds cred serets found")
  }

  username, err := vh.GetSecret("username")
  if err != nil {
    return rds, logs.Errorf("failed to get username: %v", err)
  }
  rds.User = username

  password, err := vh.GetSecret("password")
  if err != nil {
    return rds, logs.Errorf("failed to get password: %v", err)
  }
  rds.Password = password

  // Get Details
  if err := vh.GetSecrets(s.VaultDetails.DetailPath); err != nil {
    return rds, logs.Errorf("failed to get detail secrets from vault: %v", err)
  }
  if vh.Secrets() == nil {
    return rds, logs.Error("no rds detail secrets found")
  }

  port, err := vh.GetSecret("rds-port")
  if err != nil {
    fmt.Printf("\n\nporterr: %+v\n\n", err)

    if err.Error() != fmt.Sprint("key not found: rds-port") {
      return rds, logs.Errorf("failed to get port: %v", err)
    }
    port = "5432"
  }
  if port != "" {
    iport, err := strconv.Atoi(port)
    if err != nil {
      return rds, logs.Errorf("failed to parse port: %v", err)
    }
    rds.Port = iport
  }

  db, err := vh.GetSecret("rds-db")
  if err != nil {
    return rds, logs.Errorf("failed to get db: %v", err)
  }
  rds.DBName = db

  host, err := vh.GetSecret("rds-hostname")
  if err != nil {
    if err.Error() != fmt.Sprint("key not found: rds-hostname") {
      return rds, logs.Errorf("failed to get host: %v", err)
    }
    host = "db.chewed-k8s.net"
  }
  rds.Host = host

  return rds, nil
}
