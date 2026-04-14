# og-template - AI Agent Instructions

## Role

You are working on a Go CLI application built with this template. The project uses:

- **CLI**: Fang v2 + Cobra
- **Config**: Viper (YAML + env vars)
- **DI**: samber/do v2
- **Libraries**: samber/lo, samber/mo, samber/oops
- **Logging**: zerolog + slog-zerolog bridge
- **Testing**: stretchr/testify

## Commands Reference

```bash
mise exec -- task build        # Build binary
mise exec -- task test         # Run tests
mise exec -- task lint         # golangci-lint (50+ linters)
mise exec -- task fmt          # Auto-fix issues
mise exec -- task ci           # Full CI pipeline
mise exec -- task clean        # Clean artifacts
```

## Code Style

- Follow existing patterns in `internal/di/` and `cmd/myapp/`
- Use `oops.Wrapf()` for error wrapping with context
- Use `lo.Map`, `lo.SliceToMap`, `lo.MaxBy` for collections
- Use `mo.Option`, `mo.Result` for error handling
- Never ignore errors - `errcheck` is enabled
- No test exclusions - all code must pass linting

## When Adding Commands

1. Create new file in `cmd/myapp/yourcmd.go`
2. Export `newYourCmd()` function
3. Add to root command in `cmd/myapp/root.go`
4. If config needed, use existing DI services or register new ones

## When Adding Services

1. Create service in `internal/yourservice/`
2. Register in `internal/di/register.go`: `do.Provide(injector, NewYourService)`
3. Inject where needed: `svc := do.MustInvoke[*YourService](injector)`
