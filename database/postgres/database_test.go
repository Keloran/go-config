package postgres

import (
	"context"
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
	tpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"testing"
	"time"
)

func TestBuildVault(t *testing.T) {
	os.Clearenv()
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "rds-port", Value: "1111"},
			{Key: "rds-db", Value: "testDB"},
			{Key: "rds-hostname", Value: "testHost"},
		},
	}

	vd := &VaultDetails{
		CredPath:    "tester",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	db, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 1111, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}

func TestBuildVaultNoPort(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "password", Value: "testPassword"},
			{Key: "username", Value: "testUser"},
			{Key: "rds-db", Value: "testDB"},
			{Key: "rds-hostname", Value: "testHost"},
		},
	}

	vd := &VaultDetails{
		CredPath:    "tester",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	db, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 5432, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RDS_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_PORT", "1111"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_USERNAME", "testUser"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_PASSWORD", "testPassword"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("RDS_DB", "testDB"); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	db, err := d.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testPassword", db.Password)
	assert.Equal(t, "testUser", db.User)
	assert.Equal(t, 1111, db.Port)
	assert.Equal(t, "testDB", db.DBName)
	assert.Equal(t, "testHost", db.Host)
}

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

	// Start Postgres container
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
	sys.Details.ConnectionTimeout = time.Second * 30

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

	// Start Postgres container
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
	sys.Details.ConnectionTimeout = time.Second * 30

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
