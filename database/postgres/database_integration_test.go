//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tpg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPostgres(ctx context.Context) (*tpg.PostgresContainer, error) {
	pg, err := tpg.Run(ctx,
		"postgres:16-alpine",
		tpg.WithDatabase("test"),
		tpg.WithUsername("user"),
		tpg.WithPassword("password"))
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func TestPostgresPoolConnection(t *testing.T) {
	ctx := context.Background()

	pg, err := setupPostgres(ctx)
	assert.NoError(t, err)
	defer func() {
		if pg != nil {
			if err := pg.Terminate(ctx); err != nil {
				t.Logf("failed to terminate container: %v", err)
			}
		}
	}()
	assert.NotNil(t, pg)

	connectionString, err := pg.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	assert.NoError(t, err)

	sys := NewSystem()
	if err := sys.ParseConnectionString(connectionString); err != nil {
		t.Fatal(err)
	}
	sys.Details.ConnectionTimeout = 30 * time.Second

	pool, err := sys.GetPGXPoolClient(ctx)
	assert.NoError(t, err)
	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()
	assert.NotNil(t, pool)
}

func TestPostgresConnection(t *testing.T) {
	ctx := context.Background()

	pg, err := setupPostgres(ctx)
	assert.NoError(t, err)
	defer func() {
		if pg != nil {
			if err := pg.Terminate(ctx); err != nil {
				t.Logf("failed to terminate container: %v", err)
			}
		}
	}()
	assert.NotNil(t, pg)

	connectionString, err := pg.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	assert.NoError(t, err)

	sys := NewSystem()
	if err := sys.ParseConnectionString(connectionString); err != nil {
		t.Fatal(err)
	}
	sys.Details.ConnectionTimeout = 30 * time.Second

	conn, err := sys.GetPGXClient(ctx)
	assert.NoError(t, err)
	defer func() {
		if conn != nil {
			if err := conn.Close(ctx); err != nil {
				t.Logf("failed to close connection: %v", err)
			}
		}
	}()
	assert.NotNil(t, conn)
}
