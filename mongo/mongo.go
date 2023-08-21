package mongo

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vault_helper "github.com/keloran/vault-helper"
	"time"
)

type VaultHelper interface {
	GetSecrets(path string) error
	GetSecret(key string) (string, error)
	Secrets() []vault_helper.KVSecret
}

type VaultDetails struct {
	Address    string
	Path       string `env:"MONGO_VAULT_PATH" envDefault:"secret/data/chewedfeed/postgres"`
	Token      string
	ExpireTime time.Time
}

type Mongo struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost"`
	Username string `env:"MONGO_USER" envDefault:""`
	Password string `env:"MONGO_PASS" envDefault:""`
	Database string `env:"MONGO_DB" envDefault:""`

	Collections struct {
		List string `env:"MONGO_TODO_COLLECTION" envDefault:""`
	}
}

func Setup(vaultAddress, vaultToken string) VaultDetails {
	return VaultDetails{
		Address: vaultAddress,
		Token:   vaultToken,
	}
}

func Build(vd VaultDetails, vh VaultHelper) (*Mongo, error) {
	mungo := &Mongo{}

	if err := env.Parse(mungo); err != nil {
		return nil, logs.Errorf("error parsing mongo: %v", err)
	}

	v := vh.NewVault(c.Vault.Address, c.Vault.Token)
	if err := v.GetSecrets(mungo.Vault.Path); err != nil {
		return nil, logs.Errorf("error getting mongo secrets: %v", err)
	}

	username, err := v.GetSecret("username")
	if err != nil {
		return nil, logs.Errorf("error getting username: %v", err)
	}

	password, err := v.GetSecret("password")
	if err != nil {
		return nil, logs.Errorf("error getting password: %v", err)
	}

	mungo.Vault.ExpireTime = time.Now().Add(time.Duration(v.LeaseDuration) * time.Second)
	mungo.Password = password
	mungo.Username = username

	return mungo, nil
}

func GetMongoClient(ctx context.Context, cfg Config) (*mongo.Client, error) {
	if time.Now().Unix() > cfg.Mongo.Vault.ExpireTime.Unix() {
		if err := BuildMongo(&cfg); err != nil {
			return nil, logs.Errorf("error re-building mongo: %v", err)
		}
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host)))
	if err != nil {
		return nil, logs.Errorf("error connecting to mongo: %v", err)
	}

	return client, nil
}
