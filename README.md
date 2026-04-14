# og-template

Omar's Opinionated Go project template with batteries included.

## What's Included

| Category       | Tool                                                                                                              |
| -------------- | ----------------------------------------------------------------------------------------------------------------- |
| **CLI**        | [Cobra](https://github.com/spf13/cobra) + [Fang v2](https://charm.land/fang) (styled help, manpages, completions) |
| **Config**     | [Viper](https://github.com/spf13/viper) (YAML + env vars + defaults)                                              |
| **DI**         | [samber/do v2](https://github.com/samber/do) (lazy dependency injection)                                          |
| **Functional** | [samber/lo](https://github.com/samber/lo) (map, filter, reduce)                                                   |
| **Monads**     | [samber/mo](https://github.com/samber/mo) (Option, Result, Either)                                                |
| **Errors**     | [samber/oops](https://github.com/samber/oops) (structured errors with context)                                    |
| **Logging**    | [zerolog](https://github.com/rs/zerolog) + [slog-zerolog](https://github.com/samber/slog-zerolog) bridge          |
| **Testing**    | [testify](https://github.com/stretchr/testify) (assert + require)                                                 |
| **Linting**    | [golangci-lint v2](https://golangci-lint.run/) (50+ linters, strict config)                                       |
| **Tasks**      | [Task](https://taskfile.dev/) (build, test, lint, ci)                                                             |
| **Tools**      | [mise](https://mise.jdx.dev/) (Go, Task, golangci-lint, lefthook versions)                                        |
| **Hooks**      | [Lefthook](https://github.com/evilmartians/lefthook) (pre-commit, pre-push, conventional commits)                 |
| **Release**    | [GoReleaser v2](https://goreleaser.com/) (cross-compile, checksums, changelog)                                    |
| **CI/CD**      | GitHub Actions (lint + test + build matrix + release)                                                             |
| **Deps**       | [Renovate](https://docs.renovatebot.com/) (automated dependency updates)                                          |
| **AI Skills**  | [cc-skills-golang](https://github.com/samber/cc-skills-golang) (opinionated agentic coding skills in `.agents/`)  |
| **Init**       | [gum](https://github.com/charmbracelet/gum) (interactive project setup wizard via `task init`)                    |

## Quick Start

### Use this template

Click **"Use this template"** on GitHub, then:

```bash
git clone git@github.com:yourname/yourproject.git
cd yourproject
```

### Initialize your project

Run the interactive init task to rename module, binary, and env prefix:

```bash
mise install          # Install Go, Task, golangci-lint, lefthook, gum
task init             # Rename + deps + git hooks
task ci               # Verify everything works
```

`task init` uses [gum](https://github.com/charmbracelet/gum) to prompt for your module path, binary name, and env prefix, then rewrites all files, renames `cmd/myapp/`, runs `go mod tidy`, downloads deps, installs git hooks, and cleans up after itself (removes `scripts/`, gum from `.mise.toml`, and the init task from `Taskfile.yml`).

## Project Structure

```
.
├── cmd/myapp/            # CLI entrypoint and commands
│   ├── main.go           #   fang.Execute with signal handling
│   ├── root.go           #   Root cobra command
│   ├── config.go         #   config show/validate commands
│   └── version.go        #   version command
├── internal/
│   ├── config/           # Viper config loading + validation
│   │   ├── config.go     #   Config struct + Validate()
│   │   └── loader.go     #   Load() returns mo.Result[*Config]
│   ├── di/               # samber/do dependency injection
│   │   ├── container.go  #   Root container with oops errors
│   │   ├── register.go   #   Service registration
│   │   ├── config_service.go
│   │   └── logger_service.go  # zerolog + slog bridge
│   └── vinfo/            # Build version metadata (ldflags)
├── .github/workflows/
│   ├── ci.yml            # Lint + test + cross-platform build
│   └── release.yml       # GoReleaser on tag push
├── Taskfile.yml          # build, test, lint, fmt, ci, clean, setup
├── .golangci.yml         # 50+ linters, strict settings
├── .goreleaser.yaml      # Cross-compile + changelog + archives
├── .mise.toml            # Pinned tool versions
├── lefthook.yml          # Pre-commit, pre-push, conventional commits
└── config.example.yaml   # Example configuration
```

## Tasks

```bash
task              # List all tasks
task build        # Build binary with ldflags
task test         # Run tests with race detector
task test-coverage # Tests + coverage report
task lint         # golangci-lint
task fmt          # golangci-lint --fix
task ci           # fmt + lint + test + build
task clean        # Remove all artifacts and caches
task init         # Rename + deps + hooks (first-time only)
```

## Configuration

Configuration is loaded from (in order of precedence):

1. CLI flag `--config path/to/config.yaml`
2. Environment variables prefixed with `MYAPP_` (e.g. `MYAPP_APP_NAME`)
3. `config.yaml` in current directory
4. `$HOME/.config/myapp/config.yaml`
5. Built-in defaults

```bash
myapp config show       # Display resolved config
myapp config validate   # Validate config
```

## Releasing

```bash
git tag v0.1.0
git push origin v0.1.0  # Triggers GoReleaser via GitHub Actions
```

## License

MIT
