// Package di wires the application runtime dependency graph.
package di

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/samber/oops"

	"github.com/omarluq/og-template/internal/config"
)

// Container wraps the root injector used by the CLI runtime.
type Container struct {
	injector *do.RootScope
}

// NewContainer builds the root injector for the CLI runtime.
func NewContainer(configPath string) (*Container, error) {
	injector := do.New()
	do.ProvideNamedValue(injector, ConfigPathKey, configPath)
	RegisterServices(injector)

	if _, err := do.Invoke[*LoggerService](injector); err != nil {
		return nil, oops.
			In("di").
			Code("container_init").
			Wrapf(err, "initialize container")
	}

	return &Container{injector: injector}, nil
}

// Config returns the resolved application configuration.
func (c *Container) Config() *config.Config {
	return do.MustInvoke[*ConfigService](c.injector).Get()
}

// ShutdownWithContext stops all registered services using the given context.
func (c *Container) ShutdownWithContext(ctx context.Context) *do.ShutdownReport {
	return c.injector.ShutdownWithContext(ctx)
}
