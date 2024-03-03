package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type MockMongoOperations struct {
	Client     *mongo.Client
	Collection *mongo.Collection
	Database   *mongo.Database
}

func (mock *MockMongoOperations) GetMongoClient(ctx context.Context, m Mongo) (*mongo.Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if mock.Client == nil {
		return nil, errors.New("mocked error: client is nil") // Return an error when client is nil
	}

	if m.Host == "" {
		return nil, errors.New("mocked error: host is empty") // Return an error when host is empty
	}

	return mock.Client, nil
}

func (mock *MockMongoOperations) GetMongoDatabase(m Mongo) (*mongo.Database, error) {
	// return your mocked Database and error here
	if mock.Database == nil {
		return nil, errors.New("mocked error: database is nil") // Return an error when database is nil
	}

	_ = fmt.Sprintf("Mongo: %v", m)

	return mock.Database, nil
}

func (mock *MockMongoOperations) GetMongoCollection(m Mongo, collection string) (*mongo.Collection, error) {
	if collection == "" {
		return nil, errors.New("mocked error: collection is empty") // Return an error when collection is empty
	}

	_ = fmt.Sprintf("Mongo: %v", m)

	// return your mocked Collection and error here
	return mock.Collection, nil
}

func (mock *MockMongoOperations) Disconnect(ctx context.Context) error {
	return nil
}

func (mock *MockMongoOperations) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if document == nil {
		return nil, errors.New("mocked error: document is nil") // Return an error when document is nil
	}

	// return your mocked InsertOneResult and error here
	return &mongo.InsertOneResult{}, nil
}

func (mock *MockMongoOperations) InsertMany(ctx context.Context, documents []interface{}) (*mongo.InsertManyResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(documents) == 0 {
		return nil, errors.New("mocked error: documents is empty") // Return an error when documents is empty
	}

	// return your mocked InsertManyResult and error here
	return &mongo.InsertManyResult{}, nil
}

func (mock *MockMongoOperations) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil
	}

	// return your mocked SingleResult here
	return &mongo.SingleResult{}
}

func (mock *MockMongoOperations) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil, errors.New("mocked error: filter is nil") // Return an error when filter is nil
	}

	// return your mocked Cursor and error here
	return &mongo.Cursor{}, nil
}

func (mock *MockMongoOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil, errors.New("mocked error: filter is nil") // Return an error when filter is nil
	}

	if update == nil {
		return nil, errors.New("mocked error: update is nil") // Return an error when update is nil
	}

	// return your mocked UpdateResult and error here
	return &mongo.UpdateResult{}, nil
}

func (mock *MockMongoOperations) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil, errors.New("mocked error: filter is nil") // Return an error when filter is nil
	}

	if update == nil {
		return nil, errors.New("mocked error: update is nil") // Return an error when update is nil
	}

	// return your mocked UpdateResult and error here
	return &mongo.UpdateResult{}, nil
}

func (mock *MockMongoOperations) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil, errors.New("mocked error: filter is nil") // Return an error when filter is nil
	}

	// return your mocked DeleteResult and error here
	return &mongo.DeleteResult{}, nil
}

func (mock *MockMongoOperations) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if filter == nil {
		return nil, errors.New("mocked error: filter is nil") // Return an error when filter is nil
	}

	// return your mocked DeleteResult and error here
	return &mongo.DeleteResult{}, nil
}
