package context

import (
	"context"

	"github.com/onlineconf/onlineconf-go"
)

type ctxConfigModuleKey struct{}

// WithModule returns a new Context that carries value module
func WithModule(ctx context.Context, mod *onlineconf.Mod) context.Context {
	return context.WithValue(ctx, ctxConfigModuleKey{}, mod)
}

// ModuleFromContext retrieves a config module from context.
func ModuleFromContext(ctx context.Context) (*onlineconf.Mod, bool) {
	m, ok := ctx.Value(ctxConfigModuleKey{}).(*onlineconf.Mod)
	return m, ok
}

type ctxConfigModuleReloaderKey struct{}

// WithModuleReloader returns a new Context that carries value module reloader
func WithModuleReloader(ctx context.Context, mod *onlineconf.ModuleReloader) context.Context {
	return context.WithValue(ctx, ctxConfigModuleReloaderKey{}, mod)
}

// ModuleReloaderFromContext retrieves a config module from context.
func ModuleReloaderFromContext(ctx context.Context) (*onlineconf.ModuleReloader, bool) {
	m, ok := ctx.Value(ctxConfigModuleReloaderKey{}).(*onlineconf.ModuleReloader)
	return m, ok
}
