package mongo

import (
  "context"
  "errors"
  "fmt"
  "reflect"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoOperations struct {
  Client     *mongo.Client
  Collection *mongo.Collection
  Database   *mongo.Database

  FakeData map[string]interface{}
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
  if mock.FakeData == nil {
    mock.FakeData = make(map[string]interface{})
  }

  // Convert document to a map to access the 'id' field
  docMap, ok := document.(map[string]interface{})
  if !ok {
    // Attempt to use reflection if the document is not already a map
    val := reflect.ValueOf(document)
    if val.Kind() == reflect.Ptr {
      val = val.Elem()
    }
    if val.Kind() == reflect.Struct {
      docMap = make(map[string]interface{})
      for i := 0; i < val.Type().NumField(); i++ {
        field := val.Type().Field(i)
        docMap[field.Name] = val.Field(i).Interface()
      }
    } else {
      return nil, errors.New("document must be a map or a struct")
    }
  }

  id := docMap["id"].(string) // Assuming 'id' is a string and always present
  mock.FakeData[id] = document

  return &mongo.InsertOneResult{InsertedID: id}, nil
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

func (mock *MockMongoOperations) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *MockSingleResult {
  bsonFilter, ok := filter.(bson.M)
  if !ok {
    return mock.newMockSingleResult(nil, errors.New("filter must be a bson.M"))
  }

  id, ok := bsonFilter["id"].(string)
  if !ok || id == "" {
    return mock.newMockSingleResult(nil, errors.New("filter must include an 'id' field"))
  }

  document, exists := mock.FakeData[id]
  if !exists {
    return mock.newMockSingleResult(nil, mongo.ErrNoDocuments)
  }

  return mock.newMockSingleResult(document, nil)
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

type MockSingleResult struct {
  decodeFunc func(v interface{}) error
  err        error
}

// SetDecode allows setting a custom decode function for the mock SingleResult
func (msr *MockSingleResult) SetDecode(decodeFunc func(v interface{}) error) *MockSingleResult {
  msr.decodeFunc = decodeFunc
  return msr
}

// Decode calls the custom decode function if it's set, simulating mongo.SingleResult's Decode method
func (msr *MockSingleResult) Decode(v interface{}) error {
  if msr.err != nil {
    return msr.err
  }
  if msr.decodeFunc != nil {
    return msr.decodeFunc(v)
  }
  return errors.New("decode function not set")
}

// newMockSingleResult creates a new mock SingleResult with the provided document and error
func (mock *MockMongoOperations) newMockSingleResult(document interface{}, err error) *MockSingleResult {
  return &MockSingleResult{
    decodeFunc: func(v interface{}) error {
      if err != nil {
        return err
      }

      // Assuming document is already in the correct format, directly assign it
      docVal := reflect.ValueOf(document)
      if docVal.Kind() == reflect.Ptr {
        docVal = docVal.Elem()
      }

      vVal := reflect.ValueOf(v)
      if vVal.Kind() != reflect.Ptr || vVal.IsNil() {
        return errors.New("decode argument must be a non-nil pointer")
      }
      vVal.Elem().Set(docVal)

      return nil
    },
    err: err,
  }
}
