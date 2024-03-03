package mongo

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vaultHelper "github.com/keloran/vault-helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	Collections map[string]string
	Collection  string

	VaultDetails
	MongoClient MungoClient
}

type MungoClient interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
}

func NewMongo(host, username, password, database string) *Mongo {
	return &Mongo{
		Host:     host,
		Username: username,
		Password: password,
		Database: database,
	}
}

func Setup(vaultAddress, vaultToken string) VaultDetails {
	return VaultDetails{
		Address: vaultAddress,
		Token:   vaultToken,
	}
}

func Build(vd VaultDetails, vh vaultHelper.VaultHelper) (*Mongo, error) {
	mungo := NewMongo("", "", "", "")
	mungo.VaultDetails = vd

	if err := env.Parse(mungo); err != nil {
		return nil, logs.Errorf("error parsing mongo: %v", err)
	}

	// env rather than vault
	if mungo.Username != "" && mungo.Password != "" {
		return mungo, nil
	}

	if err := vh.GetSecrets(mungo.VaultDetails.Path); err != nil {
		return nil, logs.Errorf("error getting mongo secrets: %v", err)
	}

	if vh.Secrets() == nil {
		return nil, logs.Errorf("no secrets found")
	}

	username, err := vh.GetSecret("username")
	if err != nil {
		return nil, logs.Errorf("error getting username: %v", err)
	}

	password, err := vh.GetSecret("password")
	if err != nil {
		return nil, logs.Errorf("error getting password: %v", err)
	}

	mungo.VaultDetails.ExpireTime = time.Now().Add(time.Duration(vh.LeaseDuration()) * time.Second)
	mungo.Password = password
	mungo.Username = username
	mungo.Collections = BuildCollections()

	return mungo, nil
}

func BuildCollections() map[string]string {
	col := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key, val := pair[0], pair[1]
		if !strings.HasPrefix(key, "MONGO_COLLECTION_") {
			continue
		}

		colKey := strings.ToLower(strings.TrimPrefix(key, "MONGO_COLLECTION_"))
		col[colKey] = val
	}

	return col
}

// Deprecated: As of ConfigBuilder v0.5.0, use RealMongoOperations.GetMongoClient
func GetMongoClient(ctx context.Context, m Mongo) (*mongo.Client, error) {
	if time.Now().Unix() > m.VaultDetails.ExpireTime.Unix() {
		mb, err := Build(m.VaultDetails, vaultHelper.NewVault(m.VaultDetails.Address, m.VaultDetails.Token))
		if err != nil {
			return nil, logs.Errorf("error building mongo: %v", err)
		}
		m = *mb
	}

	mm := RealMongoOperations{}
	if _, err := mm.GetMongoClient(ctx, m); err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	if _, err := mm.GetMongoDatabase(m); err != nil {
		return nil, logs.Errorf("error getting mongo database: %v", err)
	}
	if _, err := mm.GetMongoCollection(m, m.Collection); err != nil {
		return nil, logs.Errorf("error getting mongo collection: %v", err)
	}

	return mm.Client, nil
}
