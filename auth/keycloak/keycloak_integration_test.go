//go:build integration
// +build integration

package keycloak

import (
	"context"
	"fmt"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	key "github.com/stillya/testcontainers-keycloak"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keycloakVersion = "24.0"
	adminUser       = "admin"
	adminPassword   = "admin"
)

func setupKeycloak(ctx context.Context) (*key.KeycloakContainer, error) {
	kc, err := key.Run(ctx, fmt.Sprintf("keycloak/keycloak:%s", keycloakVersion),
		key.WithContextPath("/auth"),
		key.WithAdminUsername(adminUser),
		key.WithAdminPassword(adminPassword),
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

	realm := gocloak.RealmRepresentation{
		Realm:   gocloak.StringP(testRealm),
		Enabled: gocloak.BoolP(true),
	}

	if _, err := client.CreateRealm(ctx, token.AccessToken, realm); err != nil {
		return fmt.Errorf("failed to create realm: %v", err)
	}

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

	kc, err := setupKeycloak(ctx)
	require.NoError(t, err)
	defer func() {
		if err := kc.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	ep, err := kc.GetAuthServerURL(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = setupTestRealm(ctx, ep)
	require.NoError(t, err)

	sys := NewSystem()
	sys.Details = Details{
		Client: testClient,
		Secret: testSecret,
		Realm:  testRealm,
		Host:   ep,
	}

	client, token, err := sys.GetClient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
}
