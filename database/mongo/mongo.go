package mongo

import (
	"context"
	"os"
	"strings"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Details struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost" vault:"mongo-hostname"`
	Username string `env:"MONGO_USER" envDefault:"" vault:"username,creds"`
	Password string `env:"MONGO_PASS" envDefault:"" vault:"password,creds"`
	Database string `env:"MONGO_DB" envDefault:"" vault:"mongo-db"`
	RawURL   string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`

	Collections map[string]string `vault:"mongo-collections"`
	Collection  string
}

type System struct {
	Context context.Context

	Details

	MongoClient MungoClient
}

type MungoOperations interface {
	GetMongoClient(m System) (*mongo.Client, error)
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
	Connect(opts ...*options.ClientOptions) (*mongo.Client, error)
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
	rab := &Details{}

	if err := env.Parse(rab); err != nil {
		return nil, logs.Errorf("mongo: unable to parse env: %v", err)
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
