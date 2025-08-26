// Package filters provides utilities to filter dynamic configuration objects.
package filters

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/zalbiraw/traefik-provider/config"
)

// HTTPRouters filters HTTP routers based on `cfg.Filter` and optional provider filter.
func HTTPRouters(routers map[string]*dynamic.Router, cfg *config.RoutersConfig, pf config.ProviderFilter) map[string]*dynamic.Router {
	result := make(map[string]*dynamic.Router)
	filters := cfg.Filter
	if pf.Provider != "" {
		filters.Provider = pf.Provider
	}

	if filters.Name == "" && filters.Provider == "" && len(filters.Entrypoints) == 0 && filters.Rule == "" && filters.Service == "" {
		return routers
	}

	filtered := filterMapByNameRegex[dynamic.Router, *dynamic.Router](routers, filters.Name, filters.Provider)
	for name, router := range filtered {
		if !httpRouterMatches(router, filters) {
			continue
		}
		result[name] = router
	}
	return result
}

func httpRouterMatches(router *dynamic.Router, filters config.RouterFilter) bool {
	if len(filters.Entrypoints) > 0 && !routerEntrypointsMatch(router.EntryPoints, filters.Entrypoints) {
		return false
	}
	if filters.Rule != "" {
		matched, err := regexMatch(filters.Rule, router.Rule)
		if err != nil || !matched {
			return false
		}
	}
	if filters.Service != "" {
		matched, err := regexMatch(filters.Service, router.Service)
		if err != nil || !matched {
			return false
		}
	}
	return true
}

// HTTPServices filters HTTP services based on `cfg.Filter` and optional provider filter.
func HTTPServices(services map[string]*dynamic.Service, cfg *config.ServicesConfig, pf config.ProviderFilter) map[string]*dynamic.Service {
	result := make(map[string]*dynamic.Service)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" {
		return services
	}

	filtered := filterMapByNameRegex[dynamic.Service, *dynamic.Service](services, filter.Name, filter.Provider)
	for name, service := range filtered {
		result[name] = service
	}
	return result
}

// HTTPMiddlewares filters HTTP middlewares based on `cfg.Filter` and optional provider filter.
func HTTPMiddlewares(middlewares map[string]*dynamic.Middleware, cfg *config.MiddlewaresConfig, pf config.ProviderFilter) map[string]*dynamic.Middleware {
	result := make(map[string]*dynamic.Middleware)
	filter := cfg.Filter
	if pf.Provider != "" {
		filter.Provider = pf.Provider
	}

	if filter.Name == "" && filter.Provider == "" {
		return middlewares
	}

	filtered := filterMapByNameRegex[dynamic.Middleware, *dynamic.Middleware](middlewares, filter.Name, filter.Provider)
	for name, middleware := range filtered {
		result[name] = middleware
	}
	return result
}
