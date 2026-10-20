package vault

import (
	"testing"

	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/assert"
)

type testDetails struct {
	Host     string            `env:"T_HOST" envDefault:"localhost" vault:"t-host"`
	Port     int               `env:"T_PORT" envDefault:"80" vault:"t-port"`
	User     string            `env:"T_USER" vault:"username,creds"`
	Name     string            `env:"T_NAME" vault:"t-name"`
	Missing  string            `env:"T_MISSING" envDefault:"keep" vault:"t-missing"`
	Cols     map[string]string `vault:"t-cols"`
	Untagged string
}

func mock(secrets ...vaultHelper.KVSecret) *vaultHelper.MockVaultHelper {
	return &vaultHelper.MockVaultHelper{KVSecrets: secrets}
}

func TestFill(t *testing.T) {
	vh := mock(
		vaultHelper.KVSecret{Key: "t-host", Value: "vaultHost"},
		vaultHelper.KVSecret{Key: "t-port", Value: "1234"},
		vaultHelper.KVSecret{Key: "username", Value: "vaultUser"},
		vaultHelper.KVSecret{Key: "t-name", Value: "vaultName"},
		vaultHelper.KVSecret{Key: "t-cols", Value: "a:b,c:d"},
	)

	t.Run("fills defaults and missing", func(t *testing.T) {
		d := testDetails{Host: "localhost", Port: 80, Missing: "keep"}
		assert.NoError(t, Fill(vh, Path{Details: "d", Credentials: "c"}, &d))
		assert.Equal(t, "vaultHost", d.Host)
		assert.Equal(t, 1234, d.Port)
		assert.Equal(t, "vaultUser", d.User)
		assert.Equal(t, "keep", d.Missing)
		assert.Equal(t, map[string]string{"a": "b", "c": "d"}, d.Cols)
	})

	t.Run("env wins", func(t *testing.T) {
		t.Setenv("T_HOST", "localhost")
		d := testDetails{Host: "localhost", Name: "set"}
		assert.NoError(t, Fill(vh, Path{Details: "d", Credentials: "c"}, &d))
		assert.Equal(t, "localhost", d.Host)
		assert.Equal(t, "set", d.Name)
	})

	t.Run("bad int", func(t *testing.T) {
		d := testDetails{}
		assert.Error(t, Fill(mock(vaultHelper.KVSecret{Key: "t-port", Value: "nope"}), Path{Credentials: "c"}, &d))
	})

	t.Run("creds field with no creds path", func(t *testing.T) {
		d := testDetails{}
		assert.Error(t, Fill(vh, Path{}, &d))
	})
}
