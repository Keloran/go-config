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
