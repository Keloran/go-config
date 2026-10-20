package keycloak

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	testRealm  = "test-realm"
	testClient = "test-client"
	testSecret = "test-secret"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("KEYCLOAK_CLIENT", testClient); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_SECRET", testSecret); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_REALM", testRealm); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	kc, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, testClient, kc.Client)
	assert.Equal(t, testSecret, kc.Secret)
	assert.Equal(t, testRealm, kc.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", kc.Host)
}
