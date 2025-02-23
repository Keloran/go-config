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
	Path        string `env:"MONGO_VAULT_PATH" envDefault:"secret/data/chewedfeed/mongo"`
	CredPath    string `env:"MONGO_VAULT_CRED_PATH" envDefault:"secret/data/chewedfeed/mongo"`
	DetailsPath string `env:"MONGO_VAULT_DETAILS_PATH" envDefault:"secret/data/chewedfeed/details"`

	ExpireTime time.Time
}

type Details struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost"`
	Username string `env:"MONGO_USER" envDefault:""`
	Password string `env:"MONGO_PASS" envDefault:""`
	Database string `env:"MONGO_DB" envDefault:""`

	Collections map[string]string
	Collection  string
}

type System struct {
	Context context.Context

	Details

	MongoClient MungoClient

	VaultDetails
	VaultHelper *vaultHelper.VaultHelper
}

type MungoOperations interface {
	GetMongoClient(ctx context.Context, m System) (*mongo.Client, error)
	GetMongoDatabase(m System) (*mongo.Database, error)
	GetMongoCollection(m System, collection string) (*mongo.Collection, error)
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Disconnect(ctx context.Context) error
}

type MungoClient interface {
	Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error)
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

func (s *System) buildVault() (*Details, error) {
	rab := &Details{}
	vh := *s.VaultHelper

	// Credentials
	if err := vh.GetSecrets(s.VaultDetails.CredPath); err != nil {
		return nil, logs.Errorf("failed to get mongo secrets: %v", err)
	}
	if vh.Secrets() == nil {
		return nil, logs.Error("no mongo secrets found")
	}

	if rab.Username == "" {
		secret, err := vh.GetSecret("username")
		if err != nil {
			return nil, logs.Errorf("failed to get username: %v", err)
		}
		rab.Username = secret
	}

	if rab.Password == "" {
		secret, err := vh.GetSecret("password")
		if err != nil {
			return nil, logs.Errorf("failed to get password: %v", err)
		}
		rab.Password = secret
	}

	// Details
	if err := vh.GetSecrets(s.VaultDetails.DetailsPath); err != nil {
		return nil, logs.Errorf("failed to get mongo details: %v", err)
	}
	if vh.Secrets() == nil {
		return nil, logs.Error("no mongo details found")
	}

	if rab.Host == "" {
		secret, err := vh.GetSecret("mongo-hostname")
		if err != nil {
			return nil, logs.Errorf("failed to get hostname: %v", err)
		}
		rab.Host = secret
	}

	if rab.Database == "" {
		secret, err := vh.GetSecret("mongo-db")
		if err != nil {
			return nil, logs.Errorf("failed to get database: %v", err)
		}
		rab.Database = secret
	}

	preCollections, err := vh.GetSecret("mongo-collections")
	if err != nil {
		return nil, logs.Errorf("failed to get collections: %v", err)
	}
	rabCollections := make(map[string]string)
	collections := strings.Split(preCollections, ",")
	for _, c := range collections {
		cols := strings.Split(c, ":")
		if len(cols) != 2 {
			return nil, logs.Errorf("collection not in correct format: %v", c)
		}
		rabCollections[cols[0]] = cols[1]
	}
	rab.Collections = rabCollections

	s.Details = *rab

	return rab, nil

}

func (s *System) buildGeneric() (*Details, error) {
	rab := &Details{}

	if err := env.Parse(rab); err != nil {
		return nil, logs.Errorf("failed to parse env: %v", err)
	}

	// Build Collections
	rab.Collections = BuildCollections()

	s.Details = *rab

	return rab, nil
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
