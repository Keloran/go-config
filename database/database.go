package database

import (
  "strconv"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
  "github.com/keloran/go-config/vault"
)

type VaultDetails struct {
	Address string
	Token   string

	Path       string `env:"RDS_VAULT_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	CredPath   string `env:"RDS_VAULT_CRED_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	DetailPath string `env:"RDS_VAULT_DETAIL_PATH" envDefault:"secret/data/chewedfeed/details"`

	ExpireTime time.Time

	Exclusive bool
}

type System struct {
	Host     string `env:"RDS_HOSTNAME" envDefault:"postgres.chewedfeed"`
	Port     int    `env:"RDS_PORT" envDefault:"5432"`
	User     string `env:"RDS_USERNAME"`
	Password string `env:"RDS_PASSWORD"`
	DBName   string `env:"RDS_DB" envDefault:"postgres"`

	VaultDetails
}

func NewDatabase(host, user, password, database string, port int) *System {
	return &System{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   database,
	}
}

func Setup(address, token string, excluive bool, paths *vault.VaultPaths) VaultDetails {
  vd := VaultDetails{
    Address: address,
    Token: token,
  }

  if paths != nil {
    vd.CredPath = paths.Database.Credentials
    vd.DetailPath = paths.Database.Details
  }

	return vd
}

func Build(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
	rds := NewDatabase("", "", "", "", 0)
	rds.VaultDetails = vd

	if vd.Exclusive {
		return vaultBuild(vd, vh)
	}

	if err := env.Parse(rds); err != nil {
		return rds, logs.Errorf("failed to parse database env: %v", err)
	}

	// env rather than vault
	if rds.User != "" && rds.Password != "" {
		return rds, logs.Error("no username or password for database")
	}

	if err := vh.GetSecrets(rds.VaultDetails.Path); err != nil {
		return rds, logs.Errorf("failed to get secret: %v", err)
	}

	if vh.Secrets() == nil {
		return rds, logs.Error("no database password found")
	}

	pass, err := vh.GetSecret("password")
	if err != nil {
		return rds, logs.Errorf("failed to get password: %v", err)
	}

	user, err := vh.GetSecret("username")
	if err != nil {
		return rds, logs.Errorf("failed to get username: %v", err)
	}

	rds.Password = pass
	rds.User = user

	return rds, nil
}

func vaultBuild(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
	rds := NewDatabase("", "", "", "", 0)
	rds.VaultDetails = vd

	// Get Credentials
	if err := vh.GetSecrets(vd.CredPath); err != nil {
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
	if err := vh.GetSecrets(vd.DetailPath); err != nil {
		return rds, logs.Errorf("failed to get detail secrets from vault: %v", err)
	}
	if vh.Secrets() == nil {
		return rds, logs.Error("no rds detail secrets found")
	}

	port, err := vh.GetSecret("rds-port")
	if err != nil {
		return rds, logs.Errorf("failed to get port: %v", err)
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
		return rds, logs.Errorf("failed to get host: %v", err)
	}
	rds.Host = host

	return rds, nil
}
