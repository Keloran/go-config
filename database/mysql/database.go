package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Details struct {
	Host     string `env:"RDS_HOSTNAME" envDefault:"mysql.chewedfeed" vault:"rds-hostname"`
	Port     int    `env:"RDS_PORT" envDefault:"3306" vault:"rds-port"`
	User     string `env:"RDS_USERNAME" vault:"username,creds"`
	Password string `env:"RDS_PASSWORD" vault:"password,creds"`
	DBName   string `env:"RDS_DB" envDefault:"chewedfeed" vault:"rds-db"`
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
		return rds, logs.Errorf("mysql: unable to parse env: %v", err)
	}

	s.Details = *rds

	return rds, nil
}

func (s *System) GetMySQLClient(ctx context.Context) (*sql.DB, error) {
	client, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", s.Details.User, s.Details.Password, s.Details.Host, s.Details.Port, s.Details.DBName))
	if err != nil {
		return nil, logs.Errorf("mysql: unable to connect: %v", err)
	}
	client.SetMaxIdleConns(10)
	client.SetMaxOpenConns(10)

	return client, nil
}

func (s *System) CloseMySQLClient(ctx context.Context, conn *sql.DB) error {
	if err := conn.Close(); err != nil {
		return logs.Errorf("mysql: unable to close connection: %v", err)
	}

	return nil
}
