package mongo

import (
  "context"
  "testing"

  "github.com/stretchr/testify/assert"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
)

func TestMockMongoOperations(t *testing.T) {
  ctx := context.Background()

  mongoOps := &MockMongoOperations{
    Client:     &mongo.Client{},     // Your mocked implementation
    Collection: &mongo.Collection{}, // Your mocked implementation
    Database:   &mongo.Database{},   // Your mocked implementation
    FakeData:   make(map[string]interface{}),
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

  t.Run("Test InsertOne", func(t *testing.T) {
    doc := map[string]interface{}{
      "id":    "123",
      "name":  "test document",
      "value": "some value",
    }
    _, err := mongoOps.InsertOne(ctx, doc)
    assert.NoError(t, err)

    // Verify document is in FakeData
    retrievedDoc, exists := mongoOps.FakeData["123"]
    assert.True(t, exists)
    assert.Equal(t, doc, retrievedDoc)
  })

  t.Run("Test FindOne", func(t *testing.T) {
    filter := bson.M{"id": "123"} // Use bson.M instead of map[string]interface{}
    result := mongoOps.FindOne(ctx, filter)
    assert.NotNil(t, result)

    var foundDoc map[string]interface{}
    err := result.Decode(&foundDoc)
    assert.NoError(t, err)
    assert.Equal(t, "test document", foundDoc["name"])
  })
}
