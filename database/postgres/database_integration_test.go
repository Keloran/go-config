//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func waitForPostgresConnection(ctx context.Context, sys *System, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		conn, err := sys.GetPGXClient(ctx)
		if err == nil {
			if closeErr := conn.Close(ctx); closeErr != nil {
				return closeErr
			}
			return nil
		}

		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}

	return lastErr
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

	err = waitForPostgresConnection(ctx, sys, 30*time.Second)
	assert.NoError(t, err)

	pool, err := sys.GetPGXPoolClient(ctx)
	assert.NoError(t, err)
	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()
	assert.NotNil(t, pool)
}

// options change the pool after go-config's defaults: a PrepareConn hook sees the caller's context
func TestPGXPoolClientOptions(t *testing.T) {
	ctx := context.Background()
	pg, err := setupPostgres(ctx)
	assert.NoError(t, err)
	defer func() {
		if pg != nil {
			_ = pg.Terminate(ctx)
		}
	}()
	connectionString, err := pg.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	sys := NewSystem()
	if err := sys.ParseConnectionString(connectionString); err != nil {
		t.Fatal(err)
	}
	sys.Details.ConnectionTimeout = 30 * time.Second
	assert.NoError(t, waitForPostgresConnection(ctx, sys, 30*time.Second))

	type who struct{}
	pool, err := sys.GetPGXPoolClient(ctx, func(c *pgxpool.Config) {
		c.MaxConns = 3
		c.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
			name, _ := ctx.Value(who{}).(string)
			_, err := conn.Exec(ctx, `SELECT set_config('app.who', $1, false)`, name)
			return err == nil, err
		}
	})
	assert.NoError(t, err)
	defer pool.Close()
	assert.Equal(t, int32(3), pool.Config().MaxConns)
	var got string
	assert.NoError(t, pool.QueryRow(context.WithValue(ctx, who{}, "ava"), `SELECT current_setting('app.who')`).Scan(&got))
	assert.Equal(t, "ava", got)
	assert.NoError(t, pool.QueryRow(context.WithValue(ctx, who{}, "sam"), `SELECT current_setting('app.who')`).Scan(&got))
	assert.Equal(t, "sam", got)
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

	err = waitForPostgresConnection(ctx, sys, 30*time.Second)
	assert.NoError(t, err)

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
