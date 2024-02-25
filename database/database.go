package database

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vault_helper "github.com/keloran/vault-helper"
	"time"
)

type VaultDetails struct {
	Address    string
	Path       string `env:"RDS_VAULT_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	Token      string
	ExpireTime time.Time
}

type Database struct {
	Host     string `env:"RDS_HOSTNAME" envDefault:"postgres.chewedfeed"`
	Port     int    `env:"RDS_PORT" envDefault:"5432"`
	User     string `env:"RDS_USERNAME"`
	Password string `env:"RDS_PASSWORD"`
	DBName   string `env:"RDS_DB" envDefault:"postgres"`

	VaultDetails
}

func NewDatabase(host, user, password, database string, port int) *Database {
	return &Database{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   database,
	}
}

func Setup(vaultAddress, vaultToken string) VaultDetails {
	return VaultDetails{
		Address: vaultAddress,
		Token:   vaultToken,
	}
}

func Build(vd VaultDetails, vh vault_helper.VaultHelper) (*Database, error) {
	rds := NewDatabase("", "", "", "", 0)
	rds.VaultDetails = vd

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
