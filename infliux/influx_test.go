package infliux

import (
  "fmt"
  vaultHelper "github.com/keloran/vault-helper"
  "github.com/stretchr/testify/assert"
  "os"
  "testing"
)

type MockVaultHelper struct {
  KVSecrets []vaultHelper.KVSecret
  Lease     int
}

func (m *MockVaultHelper) GetSecrets(path string) error {
  if path == "" {
    return fmt.Errorf("path not found: %s", path)
  }

  return nil // or simulate an error if needed
}

func (m *MockVaultHelper) GetSecret(key string) (string, error) {
  for _, s := range m.Secrets() {
    if s.Key == key {
      return s.Value, nil
    }
  }
  return "", fmt.Errorf("key not found: %s", key)
}

func (m *MockVaultHelper) Secrets() []vaultHelper.KVSecret {
  return m.KVSecrets
}

func (m *MockVaultHelper) LeaseDuration() int {
  return m.Lease
}

func TestBuildGeneric(t *testing.T) {
  os.Clearenv()

  if err := os.Setenv("INFLUX_HOSTNAME", "testHost"); err != nil {
    t.Fatal(err)
  }
  if err := os.Setenv("INFLUX_USERNAME", "testUser"); err != nil {
    t.Fatal(err)
  }
  if err := os.Setenv("INFLUX_PASSWORD", "testPassword"); err != nil {
    t.Fatal(err)
  }
  if err := os.Setenv("INFLUX_BUCKET", "testBucket"); err != nil {
    t.Fatal(err)
  }
  if err := os.Setenv("INFLUX_ORG", "testOrg"); err != nil {
    t.Fatal(err)
  }

  i := NewSystem()
  in, err := i.Build()
  assert.NoError(t, err)
  assert.Equal(t, "testPassword", in.Password)
  assert.Equal(t, "testUser", in.User)
  assert.Equal(t, "testBucket", in.Bucket)
  assert.Equal(t, "testOrg", in.Org)
  assert.Equal(t, "testHost", in.Host)
}

func TestBuildVault(t *testing.T) {
  mockVault := &MockVaultHelper{
    KVSecrets: []vaultHelper.KVSecret{
      {Key: "influx-password", Value: "testPassword"},
      {Key: "influx-username", Value: "testUser"},
      {Key: "influx-bucket", Value: "testBucket"},
      {Key: "influx-hostname", Value: "testHost"},
      {Key: "influx-org", Value: "testOrg"},
    },
  }
  
  vd := &VaultDetails{
    DetailsPath: "tester",
  }
  i := NewSystem()
  i.Setup(*vd, mockVault)
  in, err := i.Build()
  assert.NoError(t, err)
  
  assert.Equal(t, "testPassword", in.Password)
  assert.Equal(t, "testUser", in.User)
  assert.Equal(t, "testBucket", in.Bucket)
  assert.Equal(t, "testOrg", in.Org)
  assert.Equal(t, "testHost", in.Host)
}

func TestBuildVaultNoHost(t *testing.T) {
  mockVault := &MockVaultHelper{
    KVSecrets: []vaultHelper.KVSecret{
      {Key: "influx-password", Value: "testPassword"},
      {Key: "influx-username", Value: "testUser"},
      {Key: "influx-bucket", Value: "testBucket"},
      {Key: "influx-org", Value: "testOrg"},
    },
  }
  
  vd := &VaultDetails{
    DetailsPath: "tester",
  }
  i := NewSystem()
  i.Setup(*vd, mockVault)
  in, err := i.Build()
  assert.NoError(t, err)
  
  assert.Equal(t, "testPassword", in.Password)
  assert.Equal(t, "testUser", in.User)
  assert.Equal(t, "testBucket", in.Bucket)
  assert.Equal(t, "testOrg", in.Org)
  assert.Equal(t, "http://db.chewed-k8s.net:8086", in.Host)
}
