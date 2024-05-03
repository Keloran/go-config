package keycloak

import (
	"context"
  vaultHelper "github.com/keloran/vault-helper"

  "github.com/Nerzal/gocloak/v13"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type VaultDetails struct {
  Address string
  Token string

  DetailPath string `env:"KEYCLOAK_VAULT_DETAIL_PATH" envDefault:"secret/data/chewedfeed/details"`

  Exclusive bool
}

type System struct {
	Client string `env:"KEYCLOAK_CLIENT" envDefault:"" json:"client,omitempty"`
	Secret string `env:"KEYCLOAK_SECRET" envDefault:"" json:"secret,omitempty"`
	Realm  string `env:"KEYCLOAK_REALM" envDefault:"" json:"realm,omitempty"`
	Host   string `env:"KEYCLOAK_HOSTNAME" envDefault:"" json:"host,omitempty"`

  VaultDetails
}

func NewKeycloak(client, secret, realm, host string) *System {
	return &System{
		Client: client,
		Secret: secret,
		Realm:  realm,
		Host:   host,
	}
}

func Build(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
	k := NewKeycloak("", "", "", "")

  if vd.Exclusive {
    return vaultBuild(vd, vh)
  }

	if err := env.Parse(k); err != nil {
		return nil, logs.Errorf("keycloak: unable to parse keycloak: %v", err)
	}
	return k, nil
}

func Setup(address, token string, exclusive bool) VaultDetails {
  vd := VaultDetails{
    Address: address,
    Token: token,
  }

  return vd
}

func (k *System) GetClient(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	client := gocloak.NewClient(k.Host)
	token, err := client.LoginClient(ctx, k.Client, k.Secret, k.Realm)
	if err != nil {
		return nil, nil, logs.Errorf("failed to login client: %v", err)
	}

	return client, token, nil
}

func vaultBuild(vd VaultDetails, vh vaultHelper.VaultHelper) (*System, error) {
  k := NewKeycloak("", "", "", "")
  k.VaultDetails = vd

  if err := vh.GetSecrets(vd.DetailPath); err != nil {
    return k, logs.Errorf("failed to get detail secrets from vault: %v", err)
  }
  if vh.Secrets() == nil {
    return k, logs.Error("no keycloak detail secrets found")
  }

  realm, err := vh.GetSecret("keycloak-realm")
  if err != nil {
    return k, logs.Errorf("keycloak realm error: %v", err)
  }
  k.Realm = realm

  hostname, err := vh.GetSecret("keycloak-hostname")
  if err != nil {
    return k, logs.Errorf("keycloak hostname error: %v", err)
  }
  k.Host = hostname

  client, err := vh.GetSecret("keycloak-client")
  if err != nil {
    return k, logs.Errorf("keycloak client error: %v", err)
  }
  k.Client = client

  secret, err := vh.GetSecret("keycloak-secret")
  if err != nil {
    return k, logs.Errorf("keycloak secret error: %v", err)
  }
  k.Secret = secret

  return k, nil
}
