//go:build integration
// +build integration

package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func setupMongo(ctx context.Context) (*mongodb.MongoDBContainer, error) {
	mc, err := mongodb.Run(ctx, "mongo:latest")
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo: %v", err)
	}

	return mc, nil
}

func shutdownMongo(ctx context.Context, mc *mongodb.MongoDBContainer) error {
	if err := testcontainers.TerminateContainer(mc); err != nil {
		return fmt.Errorf("failed to terminate container: %v", err)
	}

	return nil
}

func TestMockMungoClient_Connect(t *testing.T) {
	os.Clearenv()
	ctx := context.Background()
	m, err := setupMongo(ctx)
	assert.NoError(t, err)
	defer func() {
		if m != nil {
			if err := shutdownMongo(ctx, m); err != nil {
				t.Logf("failed to shutdown mongo: %v", err)
			}
		}
	}()
	assert.NotNil(t, m)

	connectionString, err := m.ConnectionString(ctx)
	assert.NoError(t, err)

	err = os.Setenv("MONGO_URL", connectionString)
	assert.NoError(t, err)

	mo := NewSystem()
	_, err = mo.Build()
	assert.NoError(t, err)

	mu := RealMongoOperations{}
	conn, err := mu.GetMongoClient(*mo)
	assert.NoError(t, err)
	defer func() {
		if conn != nil {
			if err := conn.Disconnect(ctx); err != nil {
				t.Logf("failed to disconnect mongo: %v", err)
			}
		}
	}()
	assert.NotNil(t, conn)
}
