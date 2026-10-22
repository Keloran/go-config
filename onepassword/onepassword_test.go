package onepassword

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSecretReference(t *testing.T) {
	tests := []struct {
		name  string
		vault string
		key   string
		want  string
		bad   bool
	}{
		{name: "item and field", vault: "prod", key: "database.password", want: "op://prod/database/password"},
		{name: "item with dot", vault: "prod", key: "database.primary.password", want: "op://prod/database.primary/password"},
		{name: "explicit reference", vault: "ignored", key: "op://other/database/password", want: "op://other/database/password"},
		{name: "missing field", vault: "prod", key: "database.", bad: true},
		{name: "missing vault", key: "database.password", bad: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secretReference(tt.vault, tt.key)
			if tt.bad {
				if err == nil {
					t.Fatal("secretReference() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("secretReference() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("secretReference() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadNoToken(t *testing.T) {
	t.Setenv(tokenEnv, "")
	called := false
	setResolver(t, func(context.Context, string, []string) (map[string]resolution, error) {
		called = true
		return nil, nil
	})

	if err := Load("prod", map[string]string{"service.api_key": "TL515_NO_TOKEN"}); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if called {
		t.Fatal("resolver called without a token")
	}

	err := LoadRequired("prod", nil)
	if err == nil || !strings.Contains(err.Error(), tokenEnv) {
		t.Fatalf("LoadRequired() error = %v, want required-token error", err)
	}
}

func TestLoadResolvesBatchAndPreservesExistingEnvironment(t *testing.T) {
	t.Setenv(tokenEnv, "test-token")
	t.Setenv("TL515_EXISTING", "local-value")
	const secretEnv = "TL515_SECRET"
	clearEnv(t, secretEnv)

	var gotToken string
	var gotReferences []string
	setResolver(t, func(_ context.Context, token string, references []string) (map[string]resolution, error) {
		gotToken = token
		gotReferences = append([]string(nil), references...)
		return map[string]resolution{
			"op://prod/service.api_credentials/password": {value: "resolved-secret"},
		}, nil
	})

	err := Load("prod", map[string]string{
		"service.api_credentials.password": "TL515_SECRET",
		"missing.field":                    "TL515_EXISTING",
	})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if gotToken != "test-token" {
		t.Errorf("resolver token = %q, want test-token", gotToken)
	}
	wantReferences := []string{"op://prod/service.api_credentials/password"}
	if !reflect.DeepEqual(gotReferences, wantReferences) {
		t.Errorf("resolver references = %#v, want %#v", gotReferences, wantReferences)
	}
	if got := lookup(t, secretEnv); got != "resolved-secret" {
		t.Errorf("resolved environment value = %q, want resolved-secret", got)
	}
	if got := lookup(t, "TL515_EXISTING"); got != "local-value" {
		t.Errorf("existing environment value = %q, want local-value", got)
	}
}

func TestLoadUsesExplicitReferences(t *testing.T) {
	t.Setenv(tokenEnv, "test-token")
	const envName = "TL515_EXPLICIT_REFERENCE"
	clearEnv(t, envName)
	setResolver(t, func(_ context.Context, _ string, references []string) (map[string]resolution, error) {
		if !reflect.DeepEqual(references, []string{"op://another-vault/item/field"}) {
			t.Errorf("resolver references = %#v", references)
		}
		return map[string]resolution{"op://another-vault/item/field": {value: "explicit-secret"}}, nil
	})

	if err := Load("prod", map[string]string{"op://another-vault/item/field": envName}); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got := lookup(t, envName); got != "explicit-secret" {
		t.Errorf("resolved environment value = %q, want explicit-secret", got)
	}
}

func TestLoadReportsReferenceWithoutSecret(t *testing.T) {
	t.Setenv(tokenEnv, "test-token")
	setResolver(t, func(context.Context, string, []string) (map[string]resolution, error) {
		return map[string]resolution{
			"op://prod/database/password": {err: "fieldNotFound"},
		}, nil
	})

	err := Load("prod", map[string]string{"database.password": "TL515_MISSING_SECRET"})
	if err == nil {
		t.Fatal("Load() error = nil, want resolve error")
	}
	if !strings.Contains(err.Error(), "op://prod/database/password") {
		t.Errorf("error %q does not identify the failed reference", err)
	}
	if strings.Contains(err.Error(), "test-secret-value") {
		t.Errorf("error exposes a secret value: %q", err)
	}
}

func TestLoadSanitizesResolverErrors(t *testing.T) {
	t.Setenv(tokenEnv, "test-token")
	setResolver(t, func(context.Context, string, []string) (map[string]resolution, error) {
		return nil, errors.New("upstream returned secret test-secret-value")
	})

	err := Load("prod", map[string]string{"database.password": "TL515_BATCH_FAILURE"})
	if err == nil {
		t.Fatal("Load() error = nil, want resolver error")
	}
	if !strings.Contains(err.Error(), "op://prod/database/password") {
		t.Errorf("error %q does not identify the requested reference", err)
	}
	if strings.Contains(err.Error(), "test-secret-value") {
		t.Errorf("error exposes a secret value: %q", err)
	}
}

func setResolver(t *testing.T, resolver resolverFunc) {
	t.Helper()
	previous := resolveReferences
	resolveReferences = resolver
	t.Cleanup(func() { resolveReferences = previous })
}

func clearEnv(t *testing.T, name string) {
	t.Helper()
	previous, wasSet := os.LookupEnv(name)
	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("unset environment variable %s: %v", name, err)
	}
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv(name, previous)
		} else {
			_ = os.Unsetenv(name)
		}
	})
}

func lookup(t *testing.T, name string) string {
	t.Helper()
	value, ok := os.LookupEnv(name)
	if !ok {
		t.Fatalf("environment variable %s was not set", name)
	}
	return value
}
