// Package filters provides utilities to filter dynamic configuration objects.
package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefikprovider/config"
)

// tcpRouterMatchesFilter reports whether a TCP router matches the given filter.
func tcpRouterMatchesFilter(router *dynamic.TCPRouter, filter config.RouterFilter) bool {
	if len(filter.Entrypoints) > 0 {
		if !routerEntrypointsMatch(router.EntryPoints, filter.Entrypoints) {
			return false
		}
	}
	if filter.Rule != "" {
		matched, err := regexMatch(filter.Rule, router.Rule)
		if err != nil || !matched {
			return false
		}
	}
	if filter.Service != "" {
		matched, err := regexMatch(filter.Service, router.Service)
		if err != nil || !matched {
			return false
		}
	}
	return true
}

// TCPRouters filters TCP routers based on `cfg.Filter` and optional provider filter.
func TCPRouters(routers map[string]*dynamic.TCPRouter, cfg *config.RoutersConfig, pf config.ProviderFilter) map[string]*dynamic.TCPRouter {
	result := make(map[string]*dynamic.TCPRouter)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" && len(filter.Entrypoints) == 0 && filter.Rule == "" && filter.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.TCPRouter, *dynamic.TCPRouter](routers, filter.Name, filter.Provider)
	for name, router := range filtered {
		if !tcpRouterMatchesFilter(router, filter) {
			continue
		}
		result[name] = router
	}
	return result
}

// TCPServices filters TCP services based on `cfg.Filter` and optional provider filter.
func TCPServices(services map[string]*dynamic.TCPService, cfg *config.ServicesConfig, pf config.ProviderFilter) map[string]*dynamic.TCPService {
	result := make(map[string]*dynamic.TCPService)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.TCPService, *dynamic.TCPService](services, filter.Name, filter.Provider)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}

// TCPMiddlewares filters TCP middlewares based on `cfg.Filter` and optional provider filter.
func TCPMiddlewares(middlewares map[string]*dynamic.TCPMiddleware, cfg *config.MiddlewaresConfig, pf config.ProviderFilter) map[string]*dynamic.TCPMiddleware {
	result := make(map[string]*dynamic.TCPMiddleware)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" {
		return middlewares
	}

	filtered := filterMapByNameRegex[dynamic.TCPMiddleware, *dynamic.TCPMiddleware](middlewares, filter.Name, filter.Provider)
	for name, middleware := range filtered {
		result[name] = middleware
	}
	return result
}
