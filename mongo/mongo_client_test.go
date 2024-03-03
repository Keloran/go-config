package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMockMongoOperations(t *testing.T) {
	ctx := context.Background()

	mongoOps := &MockMongoOperations{
		Client:     &mongo.Client{},     // Your mocked implementation
		Collection: &mongo.Collection{}, // Your mocked implementation
		Database:   &mongo.Database{},   // Your mocked implementation
	}

	t.Run("Test GetMongoClient", func(t *testing.T) {
		client, err := mongoOps.GetMongoClient(ctx, Mongo{
			Host: "localhost",
		})
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	mongoOps.Client = nil // To simulate error
	t.Run("Test GetMongoClient Error", func(t *testing.T) {
		_, err := mongoOps.GetMongoClient(ctx, Mongo{})
		assert.Error(t, err)
	})

	t.Run("Test GetMongoDatabase", func(t *testing.T) {
		db, err := mongoOps.GetMongoDatabase(Mongo{})
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	mongoOps.Database = nil // To simulate error
	t.Run("Test GetMongoDatabase Error", func(t *testing.T) {
		_, err := mongoOps.GetMongoDatabase(Mongo{})
		assert.Error(t, err)
	})

	t.Run("Test GetMongoCollection", func(t *testing.T) {
		collection, err := mongoOps.GetMongoCollection(Mongo{}, "collection")
		assert.NoError(t, err)
		assert.NotNil(t, collection)
	})
}
