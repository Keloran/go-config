package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RealMongoOperations struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Database   *mongo.Database
}

func (r *RealMongoOperations) GetMongoClient(ctx context.Context, m System) (*mongo.Client, error) {
	if time.Now().Unix() > m.VaultDetails.ExpireTime.Unix() {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("error re-building mongo: %v", err)
		}
		m = *mr
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, m.Host)))
	if err != nil {
		return nil, logs.Errorf("error connecting to mongo: %v", err)
	}

	r.Client = client
	return client, nil
}

func (r *RealMongoOperations) GetMongoDatabase(m System) (*mongo.Database, error) {
	if time.Now().Unix() > m.VaultDetails.ExpireTime.Unix() {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("error re-building mongo: %v", err)
		}
		m = *mr
	}

	r.Database = r.Client.Database(m.Database)
	return r.Database, nil
}

func (r *RealMongoOperations) GetMongoCollection(m System, collection string) (*mongo.Collection, error) {
	if time.Now().Unix() > m.VaultDetails.ExpireTime.Unix() {
		mr := NewSystem()
		mr.Setup(m.VaultDetails, *mr.VaultHelper)
		_, err := mr.Build()
		if err != nil {
			return nil, logs.Errorf("error re-building mongo: %v", err)
		}
		m = *mr
	}

	r.Collection = r.Client.Database(m.Database).Collection(m.Collections[collection])
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
