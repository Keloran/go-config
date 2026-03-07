package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const vaultRefreshBuffer = 3600

type RealMongoOperations struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Database   *mongo.Database
}

func (r *RealMongoOperations) GetMongoClient(m System) (*mongo.Client, error) {
	if m.VaultHelper != nil && time.Now().Unix() > (m.VaultDetails.ExpireTime.Unix()-vaultRefreshBuffer) {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("mongo: unable to rebuild config: %v", err)
		}
		m = *mr
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, m.Host)
	if m.RawURL != "" {
		url = m.RawURL
	}

	client, err := mongo.Connect(options.Client().ApplyURI(url).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)), options.Client().SetReadPreference(readpref.SecondaryPreferred()))
	if err != nil {
		return nil, logs.Errorf("mongo: unable to connect: %v", err)
	}

	r.Client = client
	return client, nil
}

func (r *RealMongoOperations) GetMongoDatabase(m System) (*mongo.Database, error) {
	if m.VaultHelper != nil && time.Now().Unix() > (m.VaultDetails.ExpireTime.Unix()-vaultRefreshBuffer) {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("mongo: unable to rebuild config: %v", err)
		}
		m = *mr
	}

	if r.Client == nil {
		return nil, logs.Error("mongo: unable to get database without client")
	}

	r.Database = r.Client.Database(m.Details.Database)
	return r.Database, nil
}

func (r *RealMongoOperations) GetMongoCollection(m System, collection string) (*mongo.Collection, error) {
	if m.VaultHelper != nil && time.Now().Unix() > (m.VaultDetails.ExpireTime.Unix()-vaultRefreshBuffer) {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("mongo: unable to rebuild config: %v", err)
		}
		m = *mr
	}

	r.Collection = r.Client.Database(m.Details.Database).Collection(m.Details.Collections[collection])
	return r.Collection, nil
}

func (r *RealMongoOperations) Disconnect(ctx context.Context) error {
	return r.Client.Disconnect(ctx)
}

func (r *RealMongoOperations) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, document)
}

func (r *RealMongoOperations) InsertMany(ctx context.Context, documents []interface{}) (*mongo.InsertManyResult, error) {
	return r.Collection.InsertMany(ctx, documents)
}

func (r *RealMongoOperations) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	return r.Collection.FindOne(ctx, filter)
}

func (r *RealMongoOperations) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return r.Collection.Find(ctx, filter)
}

func (r *RealMongoOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *RealMongoOperations) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateMany(ctx, filter, update)
}

func (r *RealMongoOperations) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return r.Collection.DeleteOne(ctx, filter)
}

func (r *RealMongoOperations) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return r.Collection.DeleteMany(ctx, filter)
}
