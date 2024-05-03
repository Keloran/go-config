package infliux

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
)

type VaultDetails struct {
  Address string
  Token   string
  
	Path    string `env:"INFLUX_VAULT_PATH" envDefault:"secret/data/chewedfeed/influx"`
  CredsPath string `env:"INFLUX_VAULT_CREDS_PATH" envDefault:"secret/data/chewedfeed/influx"`
  DetailsPath string `env:"INFLUX_VAULT_DETAILS_PATH" envDefault:"secret/data/chewedfeed/details"`
  
  Exclusive bool
}

type System struct {
	Host     string `env:"INFLUX_HOSTNAME" envDefault:"http://db.chewed-k8s.net:8086"`
	User     string `env:"INFLUX_USERNAME"`
	Password string `env:"INFLUX_PASSWORD"`
	Bucket   string `env:"INFLUX_BUCKET"`
	Org      string `env:"INFLUX_ORG"`

	VaultDetails
}

func NewInflux(host, user, password, bucket, org string) *System {
	return &System{
		Host:     host,
		User:     user,
		Password: password,
		Bucket:   bucket,
		Org:      org,
	}
}

func Setup(address, token string, exclusive bool) VaultDetails {
	return VaultDetails{
		Address: address,
		Token:   token,
    
    Exclusive: exclusive,
	}
}

func vaultBuild(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
  influx := NewInflux("", "", "", "", "")
  influx.VaultDetails = vd
  
  if err := vh.GetSecrets(vd.Path); err != nil {
    return influx, logs.Errorf("failed to get influx secrets for vault: %v", err)
  }
  if vh.Secrets() == nil {
    return influx, logs.Error("no influx secrets in vault")
  }
  
  host, err := vh.GetSecret("influx-host")
  if err != nil {
    return influx, logs.Errorf("failed to get influx host: %v", err)
  }
  influx.Host = host
  
  user, err := vh.GetSecret("influx-user")
  if err != nil {
    return influx, logs.Errorf("failed to get influx user: %v", err)
  }
  influx.User = user
  
  pass, err := vh.GetSecret("influx-password")
  if err != nil {
    return influx, logs.Errorf("failed to get influx password: %v", err)
  }
  influx.Password = pass
  
  bucket, err := vh.GetSecret("influx-bucket")
  if err != nil {
    return influx, logs.Errorf("failed to get influx bucket: %v", err)
  }
  influx.Bucket = bucket
  
  org, err := vh.GetSecret("influx-org")
  if err != nil {
    return influx, logs.Errorf("failed to get influx org: %v", err)
  }
  influx.Org = org
  
  return influx, nil
}

func Build(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
	influx := NewInflux("", "", "", "", "")
	influx.VaultDetails = vd
  
  if vd.Exclusive {
    return vaultBuild(vd, vh)
  }

	if err := env.Parse(influx); err != nil {
		return influx, logs.Errorf("failed to parse influx env: %v", err)
	}

  if influx.User != "" && influx.Password != "" {
    return influx, logs.Errorf("no username of password for influx")
  }
  
  if err := vh.GetSecrets(vd.Path); err != nil {
    return influx, logs.Errorf("failed to get influx secrets for vault: %v", err)
  }
  if vh.Secrets() == nil {
    return influx, logs.Error("no influx secrets in vault")
  }
  
  host, err := vh.GetSecret("influx-host")
  if err != nil {
    return influx, logs.Errorf("failed to get influx host: %v", err)
  }
  influx.Host = host
  
  user, err := vh.GetSecret("influx-user")
  if err != nil {
    return influx, logs.Errorf("failed to get influx user: %v", err)
  }
  influx.User = user
  
  pass, err := vh.GetSecret("influx-password")
  if err != nil {
    return influx, logs.Errorf("failed to get influx password: %v", err)
  }
  influx.Password = pass
  
  bucket, err := vh.GetSecret("influx-bucket")
  if err != nil {
    return influx, logs.Errorf("failed to get influx bucket: %v", err)
  }
  influx.Bucket = bucket
  
  org, err := vh.GetSecret("influx-org")
  if err != nil {
    return influx, logs.Errorf("failed to get influx org: %v", err)
  }
  influx.Org = org
  
  return influx, nil
}
