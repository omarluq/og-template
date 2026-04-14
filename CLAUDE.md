# og-template - Project Instructions for Claude

## Project Overview

This is a Go project template with a full development toolchain:
- **CLI Framework**: Fang v2 + Cobra for command-line interfaces
- **Config**: Viper with YAML/env/defaults support
- **DI**: samber/do v2 for dependency injection
- **Libraries**: samber/lo, samber/mo, samber/oops
- **Logging**: zerolog with slog bridge (slog-zerolog)
- **Testing**: stretchr/testify
- **Tooling**: mise, Task, golangci-lint, lefthook, goreleaser

## Development Workflow

### Build & Run
```bash
mise exec -- task build    # Build binary to ./bin/
mise exec -- task run      # Build and run
mise exec -- task dev      # Run with live reload (if .air.toml exists)
```

### Testing & Quality
```bash
mise exec -- task test           # Run tests with race detector
mise exec -- task test-coverage  # Tests with HTML coverage report
mise exec -- task lint           # golangci-lint (strict: 50+ linters)
mise exec -- task fmt            # Auto-fix lint issues
mise exec -- task ci             # Full CI pipeline
```

### Project Structure
```
cmd/myapp/          # CLI entrypoint (main.go, root.go, config.go, version.go)
internal/
  config/         # Viper config loader (config.go, loader.go)
  di/             # samber/do DI container (container.go, register.go, config_service.go, logger_service.go)
  vinfo/          # Build version info (version.go)
```

## Key Patterns

### Error Handling with samber/oops
```go
import "github.com/samber/oops"

return nil, oops.
    In("config").
    Code("invalid_config").
    Wrapf(err, "load configuration")
```

### Functional Patterns with samber/lo
```go
keys := lo.Map(entries, func(e configEntry, _ int) string { return e.key })
maxLen := lo.MaxBy(keys, func(a, b string) bool { return len(a) > len(b) })
lookup := lo.SliceToMap(entries, func(e configEntry) (string, string) {
    return e.key, e.value
})
```

### Monads with samber/mo
```go
// Option pattern
env := mo.EmptyableToOption(cfg.App.Env).OrElse("development")

// Result pattern (already used in config.Load())
cfg, err := config.Load(path).Get()
```

### DI with samber/do
```go
import "github.com/samber/do/v2"

do.Provide(injector, NewConfigService)
cfg := do.MustInvoke[*ConfigService](injector)
```

## Renaming When Using Template

After cloning:
1. Replace `github.com/username/myapp` with your module path
2. Replace `myapp` binary name with your project name
3. Update `MYAPP_` env prefix in `internal/config/loader.go`
4. Rename `cmd/myapp/` to `cmd/yourproject/`

## Files to Edit When Starting a Project

1. `go.mod` - update module name
2. `cmd/myapp/main.go` - import path
3. `internal/vinfo/version.go` - import path in ldflags comment
4. `Taskfile.yml` - binary name, MAIN_PACKAGE, ldflags paths
5. `.golangci.yml` - exhaustruct include pattern
6. `.goreleaser.yaml` - project_name, binary name, owner
7. `.github/workflows/*.yml` - repo references
8. `.mise.toml` - (optional, for mise pinning)
9. `config.example.yaml` - (optional, example config)
