# Config Builder
This is the config builder for most of my projects
rather than keep writing the same code over and over again

## Build only Local and Vault
```go
cfg, err := config.Build(config.Local, config.Vault)
```

This can be used with local config services
```go
type Config struct {
  config.Config
  LocalStuffs
}

type LocalStuffs struct {
  Stuff string
}

func main() {
  cfg, err := config.Build(config.Local, config.Vault)
  if err != nil {
    panic(err)
  }

  c := Config{}
  c.Config = cfg
  c.LocalStuffs = &LocalStuffs{
    Stuff: "here"
  }
}
```

---
### This is how to create a project configuration

above is the old way, the new way is

```
import (
 ConfigBuilder "github.com/keloran/go-config"
)

type ProjectConfig struct {}

func (pc ProjectConfig) Build(cfg *ConfigBuilder) error {
  if cfg.ProjectProperties == nil {
    c.ProjectProperties = make(map[string]interface{})
  }

  c.ProjectProperties["TestProperty"] = true
  return nil
}

func main() {
  cfg, err := ConfigBuilder.Build(WithProjectConfigurator(ProjectConfig{}))
  if err != nil {
    panic(err)
  }

  fmt.Printf(cfg.ProjectProperies["TestProperty"])
}
```
