// Package onepassword loads selected 1Password fields into environment
// variables for the normal go-config subsystem builders to consume.
package onepassword

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	sdk "github.com/1password/onepassword-sdk-go"
)

const (
	tokenEnv       = "OP_SERVICE_ACCOUNT_TOKEN"
	modulePath     = "github.com/keloran/go-config/v2"
	resolveTimeout = 30 * time.Second
)

type resolution struct {
	value string
	err   string
}

type resolverFunc func(context.Context, string, []string) (map[string]resolution, error)

// resolveReferences is kept behind a function variable so option behavior can
// be tested without contacting 1Password.
var resolveReferences resolverFunc = resolveWithSDK

type mapping struct {
	secretReference string
	envName         string
}

// Load fetches the requested fields when OP_SERVICE_ACCOUNT_TOKEN is set.
// Existing environment variables are never overwritten.
func Load(vaultName string, fields map[string]string) error {
	return load(vaultName, fields, false)
}

// LoadRequired is like Load, but returns an error when the service account
// token is missing or empty.
func LoadRequired(vaultName string, fields map[string]string) error {
	return load(vaultName, fields, true)
}

func load(vaultName string, fields map[string]string, required bool) error {
	token := os.Getenv(tokenEnv)
	if token == "" {
		if required {
			return fmt.Errorf("onepassword: %s is required", tokenEnv)
		}
		return nil
	}

	pending, err := buildMappings(vaultName, fields)
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return nil
	}

	references := make([]string, 0, len(pending))
	seen := make(map[string]struct{}, len(pending))
	for _, item := range pending {
		if _, ok := seen[item.secretReference]; !ok {
			references = append(references, item.secretReference)
			seen[item.secretReference] = struct{}{}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	defer cancel()

	resolved, err := resolveReferences(ctx, token, references)
	if err != nil {
		return fmt.Errorf("onepassword: unable to resolve %s (%T)", strings.Join(references, ", "), err)
	}

	for _, item := range pending {
		result, ok := resolved[item.secretReference]
		if !ok {
			return fmt.Errorf("onepassword: no result for %s", item.secretReference)
		}
		if result.err != "" {
			return fmt.Errorf("onepassword: unable to resolve %s: %s", item.secretReference, result.err)
		}
		if err := os.Setenv(item.envName, result.value); err != nil {
			return fmt.Errorf("onepassword: unable to set environment variable %s", item.envName)
		}
	}

	return nil
}

func buildMappings(vaultName string, fields map[string]string) ([]mapping, error) {
	keys := make([]string, 0, len(fields))
	for reference := range fields {
		keys = append(keys, reference)
	}
	sort.Strings(keys)

	mappings := make([]mapping, 0, len(keys))
	for _, key := range keys {
		envName := fields[key]
		if envName == "" {
			return nil, fmt.Errorf("onepassword: environment variable name is empty for %q", key)
		}
		if _, set := os.LookupEnv(envName); set {
			continue
		}

		reference, err := secretReference(vaultName, key)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping{secretReference: reference, envName: envName})
	}
	return mappings, nil
}

func secretReference(vaultName, key string) (string, error) {
	if strings.HasPrefix(key, "op://") {
		return key, nil
	}

	separator := strings.LastIndexByte(key, '.')
	if vaultName == "" || separator <= 0 || separator == len(key)-1 {
		return "", fmt.Errorf("onepassword: invalid field mapping %q", key)
	}

	item, field := key[:separator], key[separator+1:]
	return "op://" + vaultName + "/" + item + "/" + field, nil
}

func resolveWithSDK(ctx context.Context, token string, references []string) (map[string]resolution, error) {
	client, err := sdk.NewClient(
		ctx,
		sdk.WithServiceAccountToken(token),
		sdk.WithIntegrationInfo("go-config", integrationVersion()),
	)
	if err != nil {
		return nil, err
	}

	response, err := client.Secrets().ResolveAll(ctx, references)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]resolution, len(response.IndividualResponses))
	for reference, result := range response.IndividualResponses {
		switch {
		case result.Error != nil:
			resolved[reference] = resolution{err: string(result.Error.Type)}
		case result.Content == nil:
			resolved[reference] = resolution{err: "empty response"}
		default:
			resolved[reference] = resolution{value: result.Content.Secret}
		}
	}
	return resolved, nil
}

func integrationVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if info.Main.Path == modulePath {
		return info.Main.Version
	}
	for _, dep := range info.Deps {
		if dep.Path == modulePath {
			return dep.Version
		}
	}
	return "unknown"
}
