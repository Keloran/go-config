package postgres

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
