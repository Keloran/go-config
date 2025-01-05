package keycloak

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	vaultHelper "github.com/keloran/vault-helper"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"

	key "github.com/stillya/testcontainers-keycloak"
	"github.com/stretchr/testify/assert"
)

const (
	keycloakVersion = "21.1"
	adminUser       = "admin"
	adminPassword   = "admin"
	testRealm       = "test-realm"
	testClient      = "test-client"
	testSecret      = "test-secret"
)

func TestBuildVault(t *testing.T) {
	mockVault := &vaultHelper.MockVaultHelper{
		KVSecrets: []vaultHelper.KVSecret{
			{Key: "keycloak-client", Value: "testClient"},
			{Key: "keycloak-secret", Value: "testSecret"},
			{Key: "keycloak-realm", Value: "testRealm"},
		},
	}

	vd := &VaultDetails{
		Address:     "mockAddress",
		Token:       "mockToken",
		DetailsPath: "tester",
	}
	d := NewSystem()
	d.Setup(*vd, mockVault)
	kc, err := d.Build()

	assert.NoError(t, err)

	assert.Equal(t, "testClient", kc.Client)
	assert.Equal(t, "testSecret", kc.Secret)
	assert.Equal(t, "testRealm", kc.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", kc.Host)
}

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("KEYCLOAK_CLIENT", "testClient"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_SECRET", "testSecret"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("KEYCLOAK_REALM", "testRealm"); err != nil {
		t.Fatal(err)
	}

	d := NewSystem()
	kc, err := d.Build()
	assert.NoError(t, err)

	assert.Equal(t, "testClient", kc.Client)
	assert.Equal(t, "testSecret", kc.Secret)
	assert.Equal(t, "testRealm", kc.Realm)
	assert.Equal(t, "https://keys.chewedfeed.com", kc.Host)
}

func setupKeycloak(ctx context.Context) (*key.KeycloakContainer, error) {
	kc, err := key.Run(ctx,
		"keycloak/keycloak:24.0",
		testcontainers.WithWaitStrategy(wait.ForListeningPort("8080/tcp")),
		key.WithContextPath("/auth"),
		//key.WithRealmImportFile("../testdata/realm-export.json"),
		key.WithAdminUsername("admin"),
		key.WithAdminPassword("admin"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start keycloak: %v", err)
	}

	return kc, nil
}

func setupTestRealm(ctx context.Context, uri string) error {
	client := gocloak.NewClient(uri)
	token, err := client.LoginAdmin(ctx, adminUser, adminPassword, "master")
	if err != nil {
		return fmt.Errorf("failed to login as admin: %v", err)
	}

	// Create realm
	realm := gocloak.RealmRepresentation{
		Realm:   gocloak.StringP(testRealm),
		Enabled: gocloak.BoolP(true),
	}

	if _, err := client.CreateRealm(ctx, token.AccessToken, realm); err != nil {
		return fmt.Errorf("failed to create realm: %v", err)
	}

	// Create client
	clientID := testClient
	clientSecret := testSecret
	newClient := gocloak.Client{
		ClientID:                  &clientID,
		Secret:                    &clientSecret,
		StandardFlowEnabled:       gocloak.BoolP(true),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
		ServiceAccountsEnabled:    gocloak.BoolP(true),
		Enabled:                   gocloak.BoolP(true),
	}

	_, err = client.CreateClient(ctx, token.AccessToken, testRealm, newClient)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	return nil
}

func TestKeycloakIntegration(t *testing.T) {
	ctx := context.Background()

	// Start Keycloak container
	kc, err := setupKeycloak(ctx)
	require.NoError(t, err)
	defer func() {
		if err := kc.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Setup test realm and client
	ep, err := kc.GetAuthServerURL(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = setupTestRealm(ctx, ep)
	require.NoError(t, err)

	// Test your Keycloak system
	sys := NewSystem()
	sys.Details = Details{
		Client: testClient,
		Secret: testSecret,
		Realm:  testRealm,
		Host:   ep,
	}

	// Test GetClient
	client, token, err := sys.GetClient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
}
