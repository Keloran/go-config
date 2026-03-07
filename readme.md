# Config Builder

`go-config` builds shared application configuration from environment variables, Vault, and project-specific extensions.

## Basic usage

Build only the subsystems you need:

```go
package main

import (
	"fmt"

	config "github.com/keloran/go-config"
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

	config "github.com/keloran/go-config"
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
