// Package overrides applies user-defined overrides to filtered configs.
package overrides

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
)

// stripProvider removes the provider postfix after '@' in a given name.
func stripProvider(name string) string {
	if name == "" {
		return name
	}
	if i := strings.LastIndex(name, "@"); i >= 0 {
		return name[:i]
	}
	return name
}

// StripProviderFromKeys returns a new map with keys' provider postfixes stripped.
// It is generic and can be used for routers, services, and middlewares maps.
func StripProviderFromKeys[T any](m map[string]*T) map[string]*T {
	out := make(map[string]*T, len(m))
	for name, v := range m {
		out[stripProvider(name)] = v
	}
	return out
}

// StripProviderRefsRouter strips provider postfixes from router references to service and middlewares.
// Pass middlewares as nil for router types that do not support middlewares (e.g., UDP).
func StripProviderRefsRouter(service *string, middlewares *[]string) {
	*service = stripProvider(*service)
	if middlewares != nil {
		for i := range *middlewares {
			(*middlewares)[i] = stripProvider((*middlewares)[i])
		}
	}
}

// StripProvidersHTTP keeps the same API but delegates to the minimal helpers above.
func StripProvidersHTTP(cfg *dynamic.HTTPConfiguration) {
	cfg.Routers = StripProviderFromKeys(cfg.Routers)
	cfg.Middlewares = StripProviderFromKeys(cfg.Middlewares)
	cfg.Services = StripProviderFromKeys(cfg.Services)
	for _, r := range cfg.Routers {
		StripProviderRefsRouter(&r.Service, &r.Middlewares)
	}
}

// StripProvidersTCP keeps the same API but delegates to the minimal helpers above.
func StripProvidersTCP(cfg *dynamic.TCPConfiguration) {
	cfg.Routers = StripProviderFromKeys(cfg.Routers)
	cfg.Middlewares = StripProviderFromKeys(cfg.Middlewares)
	cfg.Services = StripProviderFromKeys(cfg.Services)
	for _, r := range cfg.Routers {
		// TCP routers can have middlewares
		StripProviderRefsRouter(&r.Service, &r.Middlewares)
	}
}

// StripProvidersUDP keeps the same API but delegates to the minimal helpers above.
func StripProvidersUDP(cfg *dynamic.UDPConfiguration) {
	cfg.Routers = StripProviderFromKeys(cfg.Routers)
	cfg.Services = StripProviderFromKeys(cfg.Services)
	for _, r := range cfg.Routers {
		StripProviderRefsRouter(&r.Service, nil)
	}
}
