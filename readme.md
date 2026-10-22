# Config Builder

`go-config` builds shared application configuration from environment variables, Vault, and project-specific extensions.

## Basic usage

Build only the subsystems you need:

```go
package main

import (
	"fmt"

	config "github.com/keloran/go-config/v2"
)

func main() {
	cfg, err := config.Build(config.Local, config.Vault, config.Postgres)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Local.HTTPPort)
	fmt.Println(cfg.Database.Host)
}
```

## Project-specific configuration

Projects can extend the shared config with their own data by implementing `ProjectConfigurator`.

```go
package main

import (
	"fmt"

	config "github.com/keloran/go-config/v2"
)

type AppConfig struct {
	AppName string
	Debug   bool
}

type ProjectConfig struct{}

func (pc ProjectConfig) Build(cfg *config.Config) error {
	if cfg.ProjectProperties == nil {
		cfg.ProjectProperties = make(config.ProjectProperties)
	}

	cfg.ProjectProperties.Set("feature_x_enabled", true)
	cfg.ProjectConfig = &AppConfig{
		AppName: "example-service",
		Debug:   true,
	}

	return nil
}

func main() {
	cfg, err := config.Build(
		config.Local,
		config.WithProjectConfigurator(ProjectConfig{}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.ProjectProperties.GetValue("feature_x_enabled"))

	appCfg, ok := config.GetProjectConfig[AppConfig](cfg)
	if !ok {
		panic("missing project config")
	}

	fmt.Println(appCfg.AppName)
}
```

`GetProjectConfig[T]` expects `cfg.ProjectConfig` to be stored as `*T`.

## 1Password service account (optional)

Place `config.OnePassword` before subsystem options so it loads secrets into the environment before those subsystems parse their configuration:

```go
cfg, err := config.Build(
	config.OnePassword("prod", map[string]string{
		"postgres.password":                      "RDS_PASSWORD",
		"bugfixes_frontend_credentials.username": "BUGFIXES_AGENT_KEY",
	}),
	config.Postgres,
	config.Bugfixes,
)
```

Set `OP_SERVICE_ACCOUNT_TOKEN` for the service account. Each map key is `item.field` and each value is the environment variable to set. Item names may contain dots; the last dot separates the field. Keys beginning with `op://` are passed directly to the SDK. Existing environment variables, including values loaded from `.env`, take precedence. With no token, `OnePassword` does nothing; use `OnePasswordRequired` when startup must fail if the token is missing. A lookup failure stops config building and reports the reference without logging the secret value.

The 1Password Go SDK embeds its runtime in the application binary. Measure binary size and startup time in the consuming service before rolling it out broadly.

## Vault (optional)

Every subsystem builds from env alone. Add `config.Vault` (or use `NewConfig(vh)` / `BuildLocalVH`) to also pull secrets from Vault:

```go
cfg, err := config.Build(config.Vault, config.Redis)
```

Vault only fills fields tagged `vault:"key"` whose env var isn't set and that still hold their default; keys missing from Vault are skipped. Paths come from `cfg.VaultPaths` (details default `secret/data/chewedfeed/details`).
